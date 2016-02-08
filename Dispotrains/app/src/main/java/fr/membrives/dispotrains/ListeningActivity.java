package fr.membrives.dispotrains;

import android.accounts.Account;
import android.accounts.AccountManager;
import android.app.AlertDialog;
import android.app.ListActivity;
import android.content.ContentResolver;
import android.content.Context;
import android.content.SyncStatusObserver;
import android.os.Bundle;
import android.support.v4.widget.SwipeRefreshLayout;
import android.support.v4.widget.SwipeRefreshLayout.OnRefreshListener;
import android.util.Log;
import android.view.Menu;
import android.view.MenuInflater;
import android.view.MenuItem;
import android.view.View;

import com.google.android.gms.analytics.HitBuilders;
import com.google.android.gms.analytics.Tracker;

/**
 * Activity that listens to changes in a ContentResolver.
 */
abstract public class ListeningActivity extends ListActivity
        implements SyncStatusObserver, OnRefreshListener {
    public static final String AUTHORITY = "fr.membrives.dispotrains";
    // An account type, in the form of a domain name
    public static final String ACCOUNT_TYPE = "fr.membrives.dispotrains";
    // The account name
    public static final String ACCOUNT = "Dispotrains";
    private static final String TAG = "f.m.d.ListeningActivity";
    // Dispotrains account for synchronization.
    protected Account mAccount;
    protected Tracker mTracker;

    private Object mContentProviderHandle;

    /**
     * Create a new dummy account for the sync adapter
     *
     * @param context The application context
     */
    public static Account CreateSyncAccount(Context context) {
        // Create the account type and default account
        Account newAccount = new Account(ACCOUNT, ACCOUNT_TYPE);
        // Get an instance of the Android account manager
        AccountManager accountManager = (AccountManager) context.getSystemService(ACCOUNT_SERVICE);
        /*
         * Add the account and account type, no password or user data If successful, return the
         * Account object, otherwise report an error.
         */
        ContentResolver.setSyncAutomatically(newAccount, AUTHORITY, true);
        ContentResolver.setIsSyncable(newAccount, AUTHORITY, 1);
        ContentResolver.addPeriodicSync(newAccount, AUTHORITY, new Bundle(), 1800);
        if (accountManager.addAccountExplicitly(newAccount, null, null)) {
            Log.d(TAG, "Account added.");
        } else {
            Log.d(TAG, "Account already present.");
        }
        return newAccount;
    }

    @Override
    protected void onCreate(Bundle savedInstanceBundle) {
        super.onCreate(savedInstanceBundle);
        // Create the dummy account
        mAccount = CreateSyncAccount(this);
        DispotrainsApplication application = (DispotrainsApplication) getApplication();
        mTracker = application.getDefaultTracker();
    }

    @Override
    protected void onPause() {
        super.onPause();
        ContentResolver.removeStatusChangeListener(mContentProviderHandle);
    }

    @Override
    protected void onResume() {
        super.onResume();
        mContentProviderHandle = ContentResolver
                .addStatusChangeListener(ContentResolver.SYNC_OBSERVER_TYPE_ACTIVE, this);
        onStatusChanged(0);
    }

    public void onStatusChanged(int which) {
        final boolean isSyncing = ContentResolver.isSyncActive(mAccount, AUTHORITY);
        updateOnSync();
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                ((SwipeRefreshLayout) findViewById(R.id.swipe_refresh)).setRefreshing(isSyncing);
            }
        });
    }

    /**
     * Runs on background thread.
     */
    abstract protected void updateOnSync();

    @Override
    public void onRefresh() {
        mTracker.send(
                new HitBuilders.EventBuilder().setCategory("Action").setAction("Refresh").build());
        if (!ContentResolver.isSyncActive(mAccount, AUTHORITY)) {
            doRefresh();
        }
    }

    protected void doRefresh() {
        Bundle settingsBundle = new Bundle();
        settingsBundle.putBoolean(ContentResolver.SYNC_EXTRAS_MANUAL, true);
        settingsBundle.putBoolean(ContentResolver.SYNC_EXTRAS_EXPEDITED, true);
        ContentResolver.requestSync(mAccount, AUTHORITY, settingsBundle);
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        MenuInflater inflater = getMenuInflater();
        inflater.inflate(R.menu.main, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle item selection
        switch (item.getItemId()) {
            case R.id.action_about:
                showAbout();
            case R.id.action_refresh:
                refresh();
            default:
                return super.onOptionsItemSelected(item);
        }
    }

    private void refresh() {
        doRefresh();
    }

    private void showAbout() {
        // Inflate the about message contents
        View messageView = getLayoutInflater().inflate(R.layout.about, null, false);

        // When linking text, force to always use default color. This works
        // around a pressed color state bug.

        AlertDialog.Builder builder = new AlertDialog.Builder(this);
        builder.setIcon(R.drawable.ic_launcher);
        builder.setTitle(R.string.app_name);
        builder.setView(messageView);
        builder.setCancelable(true);
        builder.setPositiveButton("OK", null);
        builder.create();
        builder.show();
    }

}
