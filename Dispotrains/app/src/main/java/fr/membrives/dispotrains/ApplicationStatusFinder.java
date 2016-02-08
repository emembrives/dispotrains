package fr.membrives.dispotrains;

import android.content.Context;
import android.content.SharedPreferences;
import android.content.pm.PackageInfo;
import android.content.pm.PackageManager;
import android.preference.PreferenceManager;
import android.util.Log;

/**
 * Finds the application status. Inspired from http://stackoverflow.com/a/17786560.
 */
public class ApplicationStatusFinder {
    /**
     * Property name of the last version code of the app.
     */
    private static final String LAST_APP_VERSION = "last_app_version";
    /**
     * Tag used for debug messages.
     */
    private static final String TAG = "f.m.dispotrains";

    /**
     * Finds out started for the first time (ever or in the current version).<br/> <br/> Note: This
     * method is <b>not idempotent</b> only the first call will determine the proper result. Any
     * subsequent calls will only return {@link ApplicationStartStatus#NORMAL} until the app is
     * started again. So you might want to consider caching the result!
     *
     * @return the type of app start
     */
    public static ApplicationStartStatus checkAppStart(Context context) {
        PackageInfo pInfo;
        SharedPreferences sharedPreferences =
                PreferenceManager.getDefaultSharedPreferences(context);
        ApplicationStartStatus applicationStartStatus = ApplicationStartStatus.NORMAL;
        try {
            pInfo = context.getPackageManager().getPackageInfo(context.getPackageName(), 0);
            int lastVersionCode = sharedPreferences.getInt(LAST_APP_VERSION, -1);
            int currentVersionCode = pInfo.versionCode;
            applicationStartStatus = checkAppStart(currentVersionCode, lastVersionCode);
            // Update version in preferences
            sharedPreferences.edit().putInt(LAST_APP_VERSION, currentVersionCode).commit();
        } catch (PackageManager.NameNotFoundException e) {
            Log.w(TAG,
                    "Unable to determine current app version from pacakge manager.");
        }
        return applicationStartStatus;
    }

    public static ApplicationStartStatus checkAppStart(int currentVersionCode,
                                                       int lastVersionCode) {
        if (lastVersionCode == -1) {
            return ApplicationStartStatus.FIRST_TIME;
        } else if (lastVersionCode < currentVersionCode) {
            return ApplicationStartStatus.FIRST_TIME_VERSION;
        } else if (lastVersionCode > currentVersionCode) {
            Log.w(TAG, "Current version code (" + currentVersionCode +
                       ") is less then the one recognized on last startup (" +
                       lastVersionCode + ").");
            return ApplicationStartStatus.NORMAL;
        } else {
            return ApplicationStartStatus.NORMAL;
        }
    }

    /**
     * Distinguishes different kinds of app starts: <li> <ul> First start ever ({@link #FIRST_TIME})
     * </ul> <ul> First start in this version ({@link #FIRST_TIME_VERSION}) </ul> <ul> Normal app
     * start ({@link #NORMAL}) </ul>
     *
     * @author schnatterer
     */
    public enum ApplicationStartStatus {
        FIRST_TIME, FIRST_TIME_VERSION, NORMAL;
    }
}