package fr.membrives.dispotrains;

import android.content.Intent;
import android.os.Bundle;
import android.support.v4.widget.SwipeRefreshLayout;
import android.view.View;
import android.widget.ListView;

import com.google.android.gms.analytics.HitBuilders;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import fr.membrives.dispotrains.adapters.StationAdapter;
import fr.membrives.dispotrains.data.Line;
import fr.membrives.dispotrains.data.Station;
import fr.membrives.dispotrains.sync.DataSource;

public class StationListActivity extends ListeningActivity {
    private DataSource mDataSource;
    private Line mLine;
    volatile private List<Station> mStations;
    private StationAdapter mAdapter;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.station_list_activity);
        ((SwipeRefreshLayout) findViewById(R.id.swipe_refresh)).setOnRefreshListener(this);

        getActionBar().setDisplayHomeAsUpEnabled(true);
        mLine = (Line) getIntent().getExtras().getParcelable("line");
        getActionBar().setTitle(mLine.getNetwork() + " " + mLine.getId());

        mDataSource = new DataSource(this);
        mStations = new ArrayList<Station>();
        mAdapter = new StationAdapter(this, mStations);
        setListAdapter(mAdapter);
    }

    @Override
    protected void onResume() {
        super.onResume();

        mStations.clear();
        mStations.addAll(mDataSource.getStationsPerLine(mLine));
        Collections.sort(mStations);
        mAdapter.notifyDataSetChanged();

        mTracker.setScreenName("StationList~" + mLine.getNetwork() + "/" + mLine.getId());
        mTracker.send(new HitBuilders.ScreenViewBuilder().build());
    }


    @Override
    protected void updateOnSync() {
        final List<Station> stations =
                new ArrayList<Station>(mDataSource.getStationsPerLine(mLine));
        Collections.sort(stations);
        runOnUiThread(new Runnable() {
            public void run() {
                mStations.clear();
                mStations.addAll(stations);
                mAdapter.notifyDataSetChanged();
            }
        });
    }


    @Override
    protected void onListItemClick(ListView l, View v, int position, long id) {
        Station item = (Station) getListAdapter().getItem(position);
        Intent appInfo = new Intent(StationListActivity.this, StationDetailActivity.class);
        appInfo.putExtra("line", mLine);
        appInfo.putExtra("station", item.getName());
        startActivity(appInfo);
    }
}
