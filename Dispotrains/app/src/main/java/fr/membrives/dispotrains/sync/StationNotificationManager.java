package fr.membrives.dispotrains.sync;

import android.app.Notification;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.app.TaskStackBuilder;
import android.content.Context;
import android.content.Intent;
import android.graphics.BitmapFactory;
import android.support.v4.app.NotificationCompat;

import com.google.common.base.Predicate;
import com.google.common.collect.Collections2;

import java.util.Collection;

import fr.membrives.dispotrains.R;
import fr.membrives.dispotrains.StationDetailActivity;
import fr.membrives.dispotrains.data.Elevator;
import fr.membrives.dispotrains.data.Line;
import fr.membrives.dispotrains.data.Station;

/**
 * Manages notifications for non-working stations
 */
public class StationNotificationManager {
    private final Context mContext;

    public StationNotificationManager(Context context) {
        mContext = context;
    }

    /**
     * Renders a notification for a watched station whose working state changed.
     */
    public void changedWorkingState(Station station) {
        emitNotification(station, Notification.PRIORITY_DEFAULT, false);
    }

    /**
     * Renders a notification for a watched stations whose working state has not changed.
     */
    public void watchedStation(Station station) {
        if (station.getWorking()) {
            emitNotification(station, Notification.PRIORITY_MIN, false);
        } else {
            emitNotification(station, Notification.PRIORITY_LOW, true);
        }
    }

    private void emitNotification(Station station, int priority, boolean ongoing) {
        NotificationCompat.Builder mBuilder =
                new NotificationCompat.Builder(mContext).setSmallIcon(R.drawable.ic_notification)
                        .setLargeIcon(BitmapFactory
                                .decodeResource(mContext.getResources(), R.drawable.ic_launcher));
        if (station.getWorking()) {
            mBuilder.setContentTitle(station.getDisplay() + " en fonctionnement.");
            String contentText;
            if (station.getElevators().size() < 2) {
                contentText = String.format("%d ascenseur en fonctionnement.",
                        station.getElevators().size());
            } else {
                contentText = String.format("%d ascenseurs en fonctionnement.",
                        station.getElevators().size());
            }
            mBuilder.setContentText(contentText);
        } else {
            Collection<Elevator> brokenElevators =
                    Collections2.filter(station.getElevators(), new Predicate<Elevator>() {
                        @Override
                        public boolean apply(Elevator input) {
                            return !input.isWorking();
                        }
                    });
            mBuilder.setContentTitle(station.getDisplay() + " en panne");
            String contentText;
            if (brokenElevators.size() < 2) {
                contentText = String.format("%d ascenseur en panne.", brokenElevators.size());
            } else {
                contentText = String.format("%d ascenseurs en panne.", brokenElevators.size());
            }
            mBuilder.setContentText(contentText);
            mBuilder.setNumber(brokenElevators.size());

            NotificationCompat.InboxStyle inboxStyle = new NotificationCompat.InboxStyle();
            // Sets a title for the Inbox in expanded layout
            inboxStyle.setBigContentTitle(String.format("Pannes - %s:", station.getDisplay()));
            // Moves events into the expanded layout
            for (Elevator elevator : brokenElevators) {
                inboxStyle.addLine(String.format("%s: %s", elevator.getSituation(),
                        elevator.getStatusDescription()));
            }
            // Moves the expanded layout object into the notification object.
            mBuilder.setStyle(inboxStyle);
        }
        mBuilder.setWhen(station.getLastUpdate().getTime());
        mBuilder.setPriority(priority);
        mBuilder.setCategory(Notification.CATEGORY_STATUS);
        mBuilder.setOngoing(ongoing);

        Intent resultIntent = new Intent(mContext, StationDetailActivity.class);
        resultIntent.putExtra("station", station.getName());
        Line firstLine = station.getLines().iterator().next();
        resultIntent.putExtra("line", firstLine);
        TaskStackBuilder stackBuilder = TaskStackBuilder.create(mContext);
        // Adds the back stack
        stackBuilder.addParentStack(StationDetailActivity.class);
        // Adds the Intent to the top of the stack
        stackBuilder.addNextIntent(resultIntent);
        // Gets a PendingIntent containing the entire back stack
        PendingIntent resultPendingIntent =
                stackBuilder.getPendingIntent(0, PendingIntent.FLAG_UPDATE_CURRENT);

        mBuilder.setContentIntent(resultPendingIntent);
        NotificationManager mNotificationManager =
                (NotificationManager) mContext.getSystemService(Context.NOTIFICATION_SERVICE);
        // mId allows you to update the notification later on.
        mNotificationManager.notify(station.getName().hashCode(), mBuilder.build());
    }
}
