package fr.membrives.dispotrains.sync;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;
import android.util.Log;

import fr.membrives.dispotrains.data.Elevator;
import fr.membrives.dispotrains.data.Line;
import fr.membrives.dispotrains.data.Station;

/**
 * Created by etienne on 04/10/14.
 */
public class DatabaseHelper extends SQLiteOpenHelper {
    public static final String TABLE_LINES = "lines";
    public static final String TABLE_LINES_STATIONS = "line_stations";
    public static final String TABLE_STATIONS = "stations";
    public static final String TABLE_ELEVATORS = "elevators";
    public static final String COLUMN_LINE_NETWORK = "network";
    public static final String COLUMN_LINE_ID = "_id";
    public static final String COLUMN_STATION_NAME = "name";
    public static final String COLUMN_STATION_DISPLAY = "display";
    public static final String COLUMN_STATION_WORKING = "working";
    public static final String COLUMN_STATION_WATCHED = "watched";
    public static final String COLUMN_ELEVATOR_ID = "_id";
    public static final String COLUMN_ELEVATOR_SITUATION = "situation";
    public static final String COLUMN_ELEVATOR_DIRECTION = "direction";
    public static final String COLUMN_ELEVATOR_STATUS_DESC = "status_desc";
    public static final String COLUMN_ELEVATOR_STATUS_TIME = "status_time";
    public static final String COLUMN_ELEVATOR_STATION = "station";

    private static final String DATABASE_NAME = "dispotrains.db";
    private static final int DATABASE_VERSION = 2;

    // Database creation sql statement
    private static final String DATABASE_CREATE_LINE =
            "create table " + TABLE_LINES + "(" + COLUMN_LINE_ID + " TEXT primary key, " +
                    COLUMN_LINE_NETWORK + " TEXT not null);";
    private static final String DATABASE_CREATE_STATION =
            "create table " + TABLE_STATIONS + "(" + COLUMN_STATION_NAME + " TEXT primary key, " +
                    COLUMN_STATION_DISPLAY + " TEXT not null, " + COLUMN_STATION_WORKING +
                    " INTEGER not null, " + COLUMN_STATION_WATCHED + " INTEGER not null);";
    private static final String DATABASE_CREATE_LINE_STATIONS =
            "create table " + TABLE_LINES_STATIONS + "(" + COLUMN_LINE_ID + " TEXT not null, " +
                    COLUMN_STATION_NAME + " TEXT not null, FOREIGN KEY(" + COLUMN_STATION_NAME +
                    ") REFERENCES " + TABLE_STATIONS + "(" + COLUMN_STATION_NAME +
                    "), FOREIGN KEY(" + COLUMN_LINE_ID + ") REFERENCES " + TABLE_LINES + "(" +
                    COLUMN_LINE_ID + "));";
    private static final String DATABASE_CREATE_ELEVATORS =
            "create table " + TABLE_ELEVATORS + "(" + COLUMN_ELEVATOR_ID + " TEXT primary key, " +
                    COLUMN_ELEVATOR_DIRECTION + " TEXT not null, " + COLUMN_ELEVATOR_SITUATION +
                    " TEXT not null, " + COLUMN_ELEVATOR_STATUS_DESC + " TEXT, " +
                    COLUMN_ELEVATOR_STATUS_TIME + " INTEGER, " + COLUMN_ELEVATOR_STATION +
                    " TEXT not null, FOREIGN KEY(" + COLUMN_ELEVATOR_STATION + ") REFERENCES " +
                    TABLE_STATIONS + "(" + COLUMN_STATION_NAME + ") );";
    private static DatabaseHelper mInstance = null;

    private DatabaseHelper(Context context) {
        super(context, DATABASE_NAME, null, DATABASE_VERSION);
    }

    public static DatabaseHelper getInstance(Context context) {
        if (mInstance == null) {
            mInstance = new DatabaseHelper(context.getApplicationContext());
        }
        return mInstance;
    }

    @Override
    public void onCreate(SQLiteDatabase database) {
        database.execSQL(DATABASE_CREATE_LINE);
        database.execSQL(DATABASE_CREATE_STATION);
        database.execSQL(DATABASE_CREATE_LINE_STATIONS);
        database.execSQL(DATABASE_CREATE_ELEVATORS);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {
        Log.w(DatabaseHelper.class.getName(),
                "Upgrading database from version " + oldVersion + " to " + newVersion +
                        ", which will destroy all old data");
        db.execSQL("DROP TABLE IF EXISTS " + TABLE_ELEVATORS);
        db.execSQL("DROP TABLE IF EXISTS " + TABLE_LINES_STATIONS);
        db.execSQL("DROP TABLE IF EXISTS " + TABLE_STATIONS);
        db.execSQL("DROP TABLE IF EXISTS " + TABLE_LINES);
        onCreate(db);
    }

    public Cursor getAllLines() {
        Cursor cursor = this.getReadableDatabase()
                .query(TABLE_LINES, new String[]{COLUMN_LINE_NETWORK, COLUMN_LINE_ID}, null, null,
                        null, null, null);
        return cursor;
    }

    public Cursor getStations(String lineId) {
        String rawQuery = new StringBuilder().append("SELECT ").append(TABLE_STATIONS).append(".")
                .append(COLUMN_STATION_NAME).append(", ").append(TABLE_STATIONS).append(".")
                .append(COLUMN_STATION_DISPLAY).append(", ").append(TABLE_STATIONS).append(".")
                .append(COLUMN_STATION_WORKING).append(", ").append(TABLE_STATIONS).append(".")
                .append(COLUMN_STATION_WATCHED).append(" from ").append(TABLE_STATIONS)
                .append(" INNER JOIN ").append(TABLE_LINES_STATIONS).append(" ON ")
                .append(TABLE_STATIONS).append(".").append(COLUMN_STATION_NAME).append("=")
                .append(TABLE_LINES_STATIONS).append(".").append(COLUMN_STATION_NAME)
                .append(" WHERE ").append(TABLE_LINES_STATIONS).append(".").append(COLUMN_LINE_ID)
                .append("=?").toString();
        return getReadableDatabase().rawQuery(rawQuery, new String[]{lineId});
    }

    public Cursor getStation(String stationName) {
        String rawQuery = new StringBuilder().append("SELECT ").append(TABLE_STATIONS).append(".")
                .append(COLUMN_STATION_NAME).append(", ").append(TABLE_STATIONS).append(".")
                .append(COLUMN_STATION_DISPLAY).append(", ").append(TABLE_STATIONS).append(".")
                .append(COLUMN_STATION_WORKING).append(", ").append(TABLE_STATIONS).append(".")
                .append(COLUMN_STATION_WATCHED).append(" from ").append(TABLE_STATIONS)
                .append(" WHERE ").append(TABLE_STATIONS).append(".").append(COLUMN_STATION_NAME)
                .append("=?").toString();
        return getReadableDatabase().rawQuery(rawQuery, new String[]{stationName});
    }

    public Cursor getElevators(String stationName) {
        return getReadableDatabase().query(TABLE_ELEVATORS,
                new String[]{COLUMN_ELEVATOR_ID, COLUMN_ELEVATOR_SITUATION,
                        COLUMN_ELEVATOR_DIRECTION, COLUMN_ELEVATOR_STATUS_DESC,
                        COLUMN_ELEVATOR_STATUS_TIME},
                new StringBuilder().append(COLUMN_ELEVATOR_STATION).append("=?").toString(),
                new String[]{stationName}, null, null, null);
    }

    public boolean addLineToDatabase(String id, String network) {
        ContentValues values = new ContentValues();
        values.put(COLUMN_LINE_ID, id);
        values.put(COLUMN_LINE_NETWORK, network);
        return getWritableDatabase()
                .insertWithOnConflict(TABLE_LINES, null, values, SQLiteDatabase.CONFLICT_REPLACE) !=
                -1;
    }

    public boolean deleteLineFromDatabase(String id) {
        return getWritableDatabase()
                .delete(TABLE_LINES, COLUMN_LINE_ID + " =?", new String[]{id}) != 0;
    }

    public boolean addStationToDatabase(Station station) {
        ContentValues values = new ContentValues();
        values.put(COLUMN_STATION_NAME, station.getName());
        values.put(COLUMN_STATION_DISPLAY, station.getDisplay());
        values.put(COLUMN_STATION_WORKING, station.getWorking() ? 1 : 0);
        values.put(COLUMN_STATION_WATCHED, station.isWatched() ? 1 : 0);
        if (getWritableDatabase().insertWithOnConflict(TABLE_STATIONS, null, values,
                SQLiteDatabase.CONFLICT_REPLACE) == -1) {
            return false;
        }
        for (Line line : station.getLines()) {
            values = new ContentValues();
            values.put(COLUMN_STATION_NAME, station.getName());
            values.put(COLUMN_LINE_ID, line.getId());
            if (getWritableDatabase().insertWithOnConflict(TABLE_LINES_STATIONS, null, values,
                    SQLiteDatabase.CONFLICT_REPLACE) == -1) {
                return false;
            }
        }
        return true;
    }

    public boolean deleteStationFromDatabase(String name) {
        getWritableDatabase()
                .delete(TABLE_LINES_STATIONS, COLUMN_STATION_NAME + " =?", new String[]{name});
        getWritableDatabase()
                .delete(TABLE_STATIONS, COLUMN_STATION_NAME + " =?", new String[]{name});
        return true;
    }

    public boolean addElevatorToDatabase(Elevator elevator) {
        ContentValues values = new ContentValues();
        values.put(COLUMN_ELEVATOR_ID, elevator.getId());
        values.put(COLUMN_ELEVATOR_DIRECTION, elevator.getDirection());
        values.put(COLUMN_ELEVATOR_SITUATION, elevator.getSituation());
        values.put(COLUMN_ELEVATOR_STATUS_DESC, elevator.getStatusDescription());
        values.put(COLUMN_ELEVATOR_STATUS_TIME, elevator.getStatusDate().getTime());
        values.put(COLUMN_ELEVATOR_STATION, elevator.getStation().getName());
        if (getWritableDatabase().insertWithOnConflict(TABLE_ELEVATORS, null, values,
                SQLiteDatabase.CONFLICT_REPLACE) == -1) {
            return false;
        }
        return true;
    }

    public boolean deleteElevatorFromDatabase(String id) {
        return getWritableDatabase()
                .delete(TABLE_ELEVATORS, COLUMN_ELEVATOR_ID + " =?", new String[]{id}) != 0;
    }

}
