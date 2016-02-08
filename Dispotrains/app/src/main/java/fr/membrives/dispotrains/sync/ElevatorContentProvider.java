package fr.membrives.dispotrains.sync;

import android.content.ContentProvider;
import android.content.ContentResolver;
import android.content.ContentValues;
import android.content.UriMatcher;
import android.database.Cursor;
import android.net.Uri;

public class ElevatorContentProvider extends ContentProvider {
    private DatabaseHelper mHelper;

    // used for the UriMacher
    private static final int LINES = 1;
    private static final int LINE = 2;
    private static final int STATIONS = 3;
    private static final int STATION = 4;
    private static final int ELEVATORS = 5;

    private static final String AUTHORITY = "fr.membrives.dispotrains";

    private static final String BASE_PATH = "dispotrains";
    public static final Uri CONTENT_URI = Uri.parse("content://" + AUTHORITY + "/" + BASE_PATH);

    public static final String CONTENT_LINES = ContentResolver.CURSOR_DIR_BASE_TYPE + "/lines";
    public static final String CONTENT_STATIONS = ContentResolver.CURSOR_ITEM_BASE_TYPE
            + "/stations";
    public static final String CONTENT_ELEVATORS = ContentResolver.CURSOR_ITEM_BASE_TYPE
            + "/elevators";

    private static final UriMatcher sURIMatcher = new UriMatcher(UriMatcher.NO_MATCH);
    static {
        sURIMatcher.addURI(AUTHORITY, BASE_PATH + "/lines", LINES);
        sURIMatcher.addURI(AUTHORITY, BASE_PATH + "/line/*", LINE);
        sURIMatcher.addURI(AUTHORITY, BASE_PATH + "/line/stations/*", STATIONS);
        sURIMatcher.addURI(AUTHORITY, BASE_PATH + "/station/*", STATION);
        sURIMatcher.addURI(AUTHORITY, BASE_PATH + "/station/elevators/*", ELEVATORS);
    }

    @Override
    public boolean onCreate() {
        mHelper = DatabaseHelper.getInstance(getContext());
        return true;
    }

    @Override
    public Cursor query(Uri uri, String[] projection, String selection, String[] selectionArgs,
            String sortOrder) {
        Cursor cursor;
        switch (sURIMatcher.match(uri)) {
            case LINES:
                cursor = mHelper.getAllLines();
                break;
            case STATIONS:
                cursor = mHelper.getStations(uri.getLastPathSegment());
                break;
            case ELEVATORS:
                cursor = mHelper.getElevators(uri.getLastPathSegment());
                break;
            default:
                throw new IllegalArgumentException("Unknown URI: " + uri);
        }
        cursor.setNotificationUri(getContext().getContentResolver(), uri);
        return cursor;
    }

    @Override
    public String getType(Uri uri) {
        return null;
    }

    @Override
    public Uri insert(Uri uri, ContentValues values) {
        return null;
    }

    @Override
    public int delete(Uri uri, String selection, String[] selectionArgs) {
        return 0;
    }

    @Override
    public int update(Uri uri, ContentValues values, String selection, String[] selectionArgs) {
        return 0;
    }
}
