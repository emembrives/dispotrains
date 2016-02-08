package fr.membrives.dispotrains.sync;

import android.content.Context;
import android.database.Cursor;

import java.util.Date;
import java.util.HashSet;
import java.util.Set;

import fr.membrives.dispotrains.data.Elevator;
import fr.membrives.dispotrains.data.Line;
import fr.membrives.dispotrains.data.Station;

/**
 * Wrapper to the local database
 */
public class DataSource {
    private final DatabaseHelper mHelper;

    public DataSource(Context context) {
        this.mHelper = DatabaseHelper.getInstance(context);
    }

    private Line cursorToLine(Cursor cursor) {
        String network = cursor.getString(0);
        String id = cursor.getString(1);
        return new Line(id, network);
    }

    private Station cursorToStation(Cursor cursor) {
        String name = cursor.getString(0);
        String display = cursor.getString(1);
        boolean working = cursor.getInt(2) != 0;
        boolean watched = cursor.getInt(3) != 0;
        return new Station(name, display, working, watched);
    }

    private Elevator cursorToElevator(Cursor cursor) {
        String id = cursor.getString(0);
        String situation = cursor.getString(1);
        String direction = cursor.getString(2);
        String description = cursor.getString(3);
        Date time = new Date(cursor.getLong(4));
        return new Elevator(id, situation, direction, description, time);
    }

    public Set<Line> getAllLines() {
        Set<Line> lines = new HashSet<Line>();

        Cursor cursor = mHelper.getAllLines();

        cursor.moveToFirst();
        while (!cursor.isAfterLast()) {
            Line line = cursorToLine(cursor);
            lines.add(line);
            cursor.moveToNext();
        }
        // make sure to close the cursor
        cursor.close();
        return lines;
    }

    public Set<Station> getStationsPerLine(Line line) {
        Cursor cursor = mHelper.getStations(line.getId());

        Set<Station> stations = new HashSet<Station>();
        cursor.moveToFirst();
        while (!cursor.isAfterLast()) {
            Station station = cursorToStation(cursor);
            stations.add(station);
            cursor.moveToNext();
        }
        // make sure to close the cursor
        cursor.close();
        return stations;
    }

    public Set<Elevator> getElevatorsPerStation(Station station) {
        Cursor cursor = mHelper.getElevators(station.getName());

        Set<Elevator> elevators = new HashSet<Elevator>();
        cursor.moveToFirst();
        while (!cursor.isAfterLast()) {
            Elevator elevator = cursorToElevator(cursor);
            elevators.add(elevator);
            cursor.moveToNext();
        }
        // make sure to close the cursor
        cursor.close();
        return elevators;
    }

    public Station getStation(String stationName) {
        Cursor stationCursor = mHelper.getStation(stationName);
        stationCursor.moveToFirst();
        Station station = cursorToStation(stationCursor);
        stationCursor.close();
        return station;
    }

    public boolean addLineToDatabase(Line line) {
        return mHelper.addLineToDatabase(line.getId(), line.getNetwork());
    }

    public boolean deleteLineFromDatabase(Line line) {
        return mHelper.deleteLineFromDatabase(line.getId());
    }

    public boolean addStationToDatabase(Station station) {
        return mHelper.addStationToDatabase(station);
    }

    public boolean deleteStationFromDatabase(Station station) {
        return mHelper.deleteStationFromDatabase(station.getName());
    }

    public boolean addElevatorToDatabase(Elevator elevator) {
        return mHelper.addElevatorToDatabase(elevator);
    }

    public boolean deleteElevatorFromDatabase(Elevator elevator) {
        return mHelper.deleteElevatorFromDatabase(elevator.getId());
    }
}
