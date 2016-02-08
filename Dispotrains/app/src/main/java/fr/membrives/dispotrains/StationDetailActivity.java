package fr.membrives.dispotrains;

import android.app.ListActivity;
import android.content.Intent;
import android.os.Bundle;
import android.support.v4.app.NavUtils;
import android.support.v4.app.TaskStackBuilder;
import android.view.Menu;
import android.view.MenuInflater;
import android.view.MenuItem;

import com.google.android.gms.analytics.HitBuilders;
import com.google.android.gms.analytics.Tracker;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import fr.membrives.dispotrains.adapters.ElevatorAdapter;
import fr.membrives.dispotrains.data.Elevator;
import fr.membrives.dispotrains.data.Station;
import fr.membrives.dispotrains.sync.DataSource;

public class StationDetailActivity extends ListActivity {
    private Station mStation;
    private DataSource mDataSource;

    private Tracker mTracker;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.station_detail_activity);
        // Create the dummy account
        mDataSource = new DataSource(this);

        getActionBar().setDisplayHomeAsUpEnabled(true);

        DispotrainsApplication application = (DispotrainsApplication) getApplication();
        mTracker = application.getDefaultTracker();
    }

    @Override
    protected void onResume() {
        super.onResume();
        String stationName = getIntent().getExtras().getString("station");
        mStation = mDataSource.getStation(stationName);
        getActionBar().setTitle(mStation.getDisplay());

        List<Elevator> elevators =
                new ArrayList<Elevator>(mDataSource.getElevatorsPerStation(mStation));
        Collections.sort(elevators);
        setListAdapter(new ElevatorAdapter(this, elevators));

        mTracker.setScreenName("StationDetail~" + mStation.getName());
        mTracker.send(new HitBuilders.ScreenViewBuilder().build());
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        MenuInflater inflater = getMenuInflater();
        inflater.inflate(R.menu.station, menu);
        onPrepareOptionsMenu(menu);
        return true;
    }

    @Override
    public boolean onPrepareOptionsMenu(Menu menu) {
        if (mStation.isWatched()) {
            menu.findItem(R.id.action_watch).setIcon(android.R.drawable.btn_star_big_on)
                    .setTitle(R.string.action_watch_on);
        } else {
            menu.findItem(R.id.action_watch).setIcon(android.R.drawable.btn_star_big_off)
                    .setTitle(R.string.action_watch_off);
        }
        return super.onPrepareOptionsMenu(menu);
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        switch (item.getItemId()) {
            // Respond to the action bar's Up/Home button
            case android.R.id.home:
                Intent upIntent = NavUtils.getParentActivityIntent(this);
                upIntent.putExtra("line", getIntent().getExtras().getParcelable("line"));
                if (NavUtils.shouldUpRecreateTask(this, upIntent)) {
                    // This activity is NOT part of this app's task, so create a new task
                    // when navigating up, with a synthesized back stack.
                    TaskStackBuilder.create(this)
                            // Add all of this activity's parents to the back stack
                            .addNextIntentWithParentStack(upIntent)
                                    // Navigate up to the closest parent
                            .startActivities();
                } else {
                    // This activity is part of this app's task, so simply
                    // navigate up to the logical parent activity.
                    NavUtils.navigateUpTo(this, upIntent);
                }
                return true;
            case R.id.action_watch:
                toggleWatch();
                return true;
        }
        return super.onOptionsItemSelected(item);
    }

    private void toggleWatch() {
        mStation.setWatched(!mStation.isWatched());
        mDataSource.addStationToDatabase(mStation);
        invalidateOptionsMenu();
    }
}
