package fr.membrives.dispotrains.sync;

import android.accounts.Account;
import android.content.AbstractThreadedSyncAdapter;
import android.content.ContentProviderClient;
import android.content.Context;
import android.content.SyncResult;
import android.os.Bundle;
import android.util.Log;

import com.google.common.collect.HashMultimap;
import com.google.common.collect.Multimap;

import org.apache.commons.io.IOUtils;
import org.joda.time.format.ISODateTimeFormat;
import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.IOException;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.net.URL;
import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import fr.membrives.dispotrains.data.Elevator;
import fr.membrives.dispotrains.data.Line;
import fr.membrives.dispotrains.data.Station;

public class SyncAdapter extends AbstractThreadedSyncAdapter {
    private static final String TAG = "f.m.d.SyncAdapter";
    // Global variables
    // Define a variable to contain a content resolver instance
    private final DataSource mSource;
    private final StationNotificationManager mNotificationManager;

    /**
     * Set up the sync adapter
     */
    public SyncAdapter(Context context, boolean autoInitialize) {
        super(context, autoInitialize);
        /*
         * If your app uses a content resolver, get an instance of it from the incoming Context
         */
        mSource = new DataSource(context);
        mNotificationManager = new StationNotificationManager(context);
    }

    /**
     * Set up the sync adapter. This form of the constructor maintains compatibility with Android
     * 3.0 and later platform versions
     */
    public SyncAdapter(Context context, boolean autoInitialize, boolean allowParallelSyncs) {
        super(context, autoInitialize, allowParallelSyncs);
        /*
         * If your app uses a content resolver, get an instance of it from the incoming Context
         */
        mSource = new DataSource(context);
        mNotificationManager = new StationNotificationManager(context);
    }

    @Override
    public void onPerformSync(Account account, Bundle extras, String authority,
                              ContentProviderClient provider, SyncResult syncResult) {
        Multimap<Line, Station> lines = HashMultimap.create();
        try {
            URL stationsURL = new URL("http://dispotrains.membrives.fr/app/GetStations/");
            InputStream stationsData = stationsURL.openStream();
            JSONArray jsonStationList = new JSONArray(IOUtils.toString(stationsData));
            for (int i = 0; i < jsonStationList.length(); i++) {
                JSONObject jsonStation = jsonStationList.getJSONObject(i);
                Multimap<Line, Station> station = processStation(jsonStation);
                lines.putAll(station);
            }
            for (Map.Entry<Line, Station> entry : lines.entries()) {
                entry.getValue().addToLine(entry.getKey());
            }
        } catch (MalformedURLException e) {
            Log.e(TAG, "Unable to sync", e);
            return;
        } catch (IOException e) {
            Log.e(TAG, "Unable to sync", e);
            return;
        } catch (JSONException e) {
            Log.e(TAG, "Unable to sync", e);
            return;
        }

        synchronizeWithDatabase(lines.keySet(), syncResult);
    }

    private void synchronizeWithDatabase(Set<Line> lines, SyncResult syncResult) {
        Set<Line> oldLines = mSource.getAllLines();
        int nbLines = 0, nbStations = 0, nbElevators = 0;
        for (Line line : lines) {
            nbLines++;
            Map<Station, Station> oldStations = new HashMap<Station, Station>();
            if (!oldLines.contains(line)) {
                mSource.addLineToDatabase(line);
                syncResult.stats.numInserts++;
                syncResult.stats.numEntries++;
            } else {
                for (Station station : mSource.getStationsPerLine(line)) {
                    oldStations.put(station, station);
                }
                oldLines.remove(line);
            }
            for (Station station : line.getStations()) {
                nbStations++;
                Map<Elevator, Elevator> oldElevators = new HashMap<Elevator, Elevator>();
                if (!oldStations.containsKey(station)) {
                    mSource.addStationToDatabase(station);
                    syncResult.stats.numInserts++;
                    syncResult.stats.numEntries++;
                } else {
                    Station oldStation = oldStations.get(station);
                    station.setWatched(oldStation.isWatched());
                    if (station.getWorking() != oldStation.getWorking()) {
                        mSource.addStationToDatabase(station);
                        if (station.isWatched()) {
                            mNotificationManager.changedWorkingState(station);
                        }
                        syncResult.stats.numUpdates++;
                        syncResult.stats.numEntries++;
                    } else if (station.isWatched()) {
                        mNotificationManager.watchedStation(station);
                    }
                    for (Elevator oldElevator : mSource.getElevatorsPerStation(station)) {
                        oldElevators.put(oldElevator, oldElevator);
                    }
                    oldStations.remove(station);
                }
                for (Elevator elevator : station.getElevators()) {
                    nbElevators++;
                    if (!oldElevators.containsKey(elevator)) {
                        mSource.addElevatorToDatabase(elevator);
                        syncResult.stats.numInserts++;
                        syncResult.stats.numEntries++;
                    } else {
                        if (!elevator.getStatusDate()
                                .equals(oldElevators.get(elevator).getStatusDate())) {
                            // Elevator status has been updated
                            mSource.addElevatorToDatabase(elevator);
                            syncResult.stats.numUpdates++;
                            syncResult.stats.numEntries++;
                        }
                        oldElevators.remove(elevator);
                    }
                }
                for (Elevator elevator : oldElevators.keySet()) {
                    mSource.deleteElevatorFromDatabase(elevator);
                    syncResult.stats.numDeletes++;
                    syncResult.stats.numEntries++;
                }
            }
            for (Station station : oldStations.keySet()) {
                mSource.deleteStationFromDatabase(station);
                syncResult.stats.numDeletes++;
                syncResult.stats.numEntries++;
            }
        }
        for (Line line : oldLines) {
            mSource.deleteLineFromDatabase(line);
            syncResult.stats.numDeletes++;
            syncResult.stats.numEntries++;
        }
        Log.d(TAG, new StringBuilder("lines, stations, elevators: ").append(nbLines).append(" ")
                .append(nbStations).append(" ").append(nbElevators).toString());
        Log.d(TAG, new StringBuilder("SyncStats: ").append(syncResult.stats.numEntries).append(" ")
                .append(syncResult.stats.numInserts).append(" ").append(syncResult.stats.numUpdates)
                .append(" ").append(syncResult.stats.numDeletes).toString());
    }

    private Multimap<Line, Station> processStation(JSONObject jsonStation) throws JSONException {
        String name = jsonStation.getString("name");
        String displayName = jsonStation.getString("displayname");

        List<Elevator> elevators = new ArrayList<Elevator>();
        boolean stationWorking = true;
        JSONArray jsonElevators = jsonStation.getJSONArray("elevators");
        for (int i = 0; i < jsonElevators.length(); i++) {
            JSONObject jsonElevator = jsonElevators.getJSONObject(i);
            JSONObject jsonStatus = jsonElevator.getJSONObject("status");
            Date lastUpdate = ISODateTimeFormat.dateTimeParser()
                    .parseDateTime(jsonStatus.getString("lastupdate")).toDate();
            Elevator elevator =
                    new Elevator(jsonElevator.getString("id"), jsonElevator.getString("situation"),
                            jsonElevator.getString("direction"), jsonStatus.getString("state"),
                            lastUpdate);
            elevators.add(elevator);
            if (stationWorking && !elevator.isWorking()) {
                stationWorking = false;
            }
        }
        Station station = new Station(name, displayName, stationWorking, false);
        for (Elevator elevator : elevators) {
            station.addElevator(elevator);
        }

        Multimap<Line, Station> lines = HashMultimap.create();
        JSONArray jsonLines = jsonStation.getJSONArray("lines");
        for (int i = 0; i < jsonLines.length(); i++) {
            JSONObject jsonLine = jsonLines.getJSONObject(i);
            Line line = new Line(jsonLine.getString("id"), jsonLine.getString("network"));
            lines.put(line, station);
        }
        return lines;
    }
}
