package fr.membrives.dispotrains.data;

import java.util.HashSet;
import java.util.Set;

import android.os.Parcel;
import android.os.Parcelable;

/**
 * A line
 */
public class Line implements Comparable<Line>, Parcelable {
    private final String id;
    private final String network;
    private final Set<Station> stations;

    public Line(String id, String network) {
        this.id = id;
        this.network = network;
        this.stations = new HashSet<Station>();
    }

    public String getNetwork() {
        return network;
    }

    public String getId() {
        return id;
    }

    public void addStation(Station station) {
        stations.add(station);
    }

    public Set<Station> getStations() {
        return stations;
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int result = 1;
        result = prime * result + ((id == null) ? 0 : id.hashCode());
        result = prime * result + ((network == null) ? 0 : network.hashCode());
        return result;
    }

    @Override
    public boolean equals(Object obj) {
        if (this == obj)
            return true;
        if (obj == null)
            return false;
        if (getClass() != obj.getClass())
            return false;
        Line other = (Line) obj;
        if (id == null) {
            if (other.id != null)
                return false;
        } else if (!id.equals(other.id))
            return false;
        if (network == null) {
            if (other.network != null)
                return false;
        } else if (!network.equals(other.network))
            return false;
        return true;
    }

    public int compareTo(Line another) {
        String lhsString = getNetwork() + " " + getId();
        String rhsString = another.getNetwork() + " " + another.getId();
        return lhsString.compareTo(rhsString);
    }

    public int describeContents() {
        // TODO Auto-generated method stub
        return 0;
    }

    public void writeToParcel(Parcel dest, int flags) {
        dest.writeString(id);
        dest.writeString(network);
    }

    public static final Parcelable.Creator<Line> CREATOR = new Parcelable.Creator<Line>() {
        public Line createFromParcel(Parcel source) {
            String id = source.readString();
            String network = source.readString();
            return new Line(id, network);
        }

        public Line[] newArray(int size) {
            return new Line[size];
        }
    };
}
