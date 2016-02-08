package fr.membrives.dispotrains;

import android.content.Intent;
import android.os.AsyncTask;
import android.os.Bundle;
import android.support.v4.widget.SwipeRefreshLayout;
import android.view.View;
import android.widget.ListView;

import com.github.amlcurran.showcaseview.OnShowcaseEventListener;
import com.github.amlcurran.showcaseview.ShowcaseView;
import com.google.android.gms.analytics.HitBuilders;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import fr.membrives.dispotrains.adapters.LineAdapter;
import fr.membrives.dispotrains.data.Line;
import fr.membrives.dispotrains.sync.DataSource;

public class MainActivity extends ListeningActivity {

    private DataSource mDataSource;
    volatile private List<Line> mLines;
    private LineAdapter mAdapter;
    private boolean mShowcaseShown = false;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        // Turn on automatic syncing for the default account and authority
        ((SwipeRefreshLayout) findViewById(R.id.swipe_refresh)).setOnRefreshListener(this);

        mDataSource = new DataSource(this);

        mLines = new ArrayList<Line>(mDataSource.getAllLines());
        Collections.sort(mLines);
        mAdapter = new LineAdapter(this, mLines);
        setListAdapter(mAdapter);

        doRefresh();
        onStatusChanged(0);

        switch (ApplicationStatusFinder.checkAppStart(MainActivity.this)) {
            case FIRST_TIME:
                new ShowcaseView.Builder(MainActivity.this).withMaterialShowcase()
                        .setStyle(R.style.ShowcaseView_Light)
                        .setContentTitle(R.string.tutorial_title)
                        .setContentText(R.string.tutorial_content)
                        .hideOnTouchOutside()
                        .setShowcaseEventListener(new OnShowcaseEventListener() {
                            @Override
                            public void onShowcaseViewHide(ShowcaseView showcaseView) {}

                            @Override
                            public void onShowcaseViewDidHide(ShowcaseView showcaseView) {
                                mShowcaseShown = false;
                                AsyncTask.execute(new Runnable() {
                                    @Override
                                    public void run() {
                                        MainActivity.this.updateOnSync();

                                    }
                                });
                            }

                            @Override
                            public void onShowcaseViewShow(ShowcaseView showcaseView) {
                                mShowcaseShown = true;
                            }
                        }).build();
                break;
            case FIRST_TIME_VERSION:
                break;
            default:
                break;
        }
    }

    @Override
    protected void onResume() {
        super.onResume();
        mTracker.setScreenName("LineList");
        mTracker.send(new HitBuilders.ScreenViewBuilder().build());
    }

    @Override
    public void updateOnSync() {
        final List<Line> lines = new ArrayList<Line>(mDataSource.getAllLines());
        Collections.sort(lines);
        if (!mShowcaseShown) {
            // Do not display anything while we are showing the tutorial.
            runOnUiThread(new Runnable() {
                public void run() {
                    mLines.clear();
                    mLines.addAll(lines);
                    mAdapter.notifyDataSetChanged();
                }
            });
        }
    }

    @Override
    protected void onListItemClick(ListView l, View v, int position, long id) {
        Line item = (Line) getListAdapter().getItem(position);
        Intent appInfo = new Intent(MainActivity.this, StationListActivity.class);
        appInfo.putExtra("line", item);
        startActivity(appInfo);
    }
}
