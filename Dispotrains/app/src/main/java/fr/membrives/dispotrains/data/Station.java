package fr.membrives.dispotrains.data;

import android.os.Parcel;
import android.os.Parcelable;

import java.util.Date;
import java.util.HashSet;
import java.util.Set;

/**
 * A train station
 */
public class Station implements Comparable<Station>, Parcelable {
    public static final Parcelable.Creator<Station> CREATOR = new Parcelable.Creator<Station>() {
        public Station createFromParcel(Parcel source) {
            String name = source.readString();
            String display = source.readString();
            // boolean[]{working, watched}
            boolean[] booleanArray = new boolean[2];
            source.readBooleanArray(booleanArray);
            return new Station(name, display, booleanArray[0], booleanArray[1]);
        }

        public Station[] newArray(int size) {
            return new Station[size];
        }
    };
    private final String name;
    private final String display;
    private final boolean working;
    private final Set<Line> lines;
    private final Set<Elevator> elevators;
    private boolean watched;

    public Station(String name, String display, boolean working, boolean watched) {
        this.name = name;
        this.display = display;
        this.working = working;
        this.watched = watched;
        this.lines = new HashSet<Line>();
        this.elevators = new HashSet<Elevator>();
    }

    public String getName() {
        return name;
    }

    public String getDisplay() {
        return display;
    }

    public boolean getWorking() {
        return working;
    }

    public void addToLine(Line line) {
        lines.add(line);
        line.addStation(this);
    }

    public void addElevator(Elevator elevator) {
        elevators.add(elevator);
        elevator.setStation(this);
    }

    public Set<Elevator> getElevators() {
        return elevators;
    }

    public boolean isWatched() {
        return watched;
    }

    public void setWatched(boolean watched) {
        this.watched = watched;
    }

    public Date getLastUpdate() {
        if (elevators.isEmpty()) {
            return new Date();
        }
        Date time = null;
        for (Elevator elevator : elevators) {
            if (time == null || elevator.getStatusDate().after(time)) {
                time = elevator.getStatusDate();
            }
        }
        return time;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int result = 1;
        result = prime * result + ((display == null) ? 0 : display.hashCode());
        result = prime * result + ((name == null) ? 0 : name.hashCode());
        result = prime * result + (working ? 1231 : 1237);
        return result;
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj) {
            return true;
        }
        if (obj == null) {
            return false;
        }
        if (getClass() != obj.getClass()) {
            return false;
        }
        Station other = (Station) obj;
        if (display == null) {
            if (other.display != null) {
                return false;
            }
        } else if (!display.equals(other.display)) {
            return false;
        }
        if (name == null) {
            if (other.name != null) {
                return false;
            }
        } else if (!name.equals(other.name)) {
            return false;
        }
        if (working != other.working) {
            return false;
        }
        return true;
    }

    public Set<Line> getLines() {
        return lines;
    }

    public int compareTo(Station another) {
        if (this.working == another.working) {
            return getDisplay().compareTo(another.getDisplay());
        } else {
            return this.working ? 1 : -1;
        }
    }

    public int describeContents() {
        return 0;
    }

    public void writeToParcel(Parcel dest, int flags) {
        dest.writeString(name);
        dest.writeString(display);
        dest.writeBooleanArray(new boolean[]{working, watched});
    }
}
