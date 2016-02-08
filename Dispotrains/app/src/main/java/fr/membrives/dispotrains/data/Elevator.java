package fr.membrives.dispotrains.data;

import java.util.Date;

/**
 * Created by etienne on 04/10/14.
 */
public class Elevator implements Comparable<Elevator> {
    private final String id;
    private final String situation;
    private final String direction;
    private final String statusDescription;
    private final Date statusDate;
    private Station station;

    public Elevator(String id, String situation, String direction, String statusDescription,
                    Date statusDate) {
        this.id = id;
        this.situation = situation;
        this.direction = direction;
        this.statusDescription = statusDescription;
        this.statusDate = statusDate;
    }

    public String getId() {
        return id;
    }

    public String getSituation() {
        return situation;
    }

    public String getDirection() {
        return direction;
    }

    public String getStatusDescription() {
        return statusDescription;
    }

    public Date getStatusDate() {
        return statusDate;
    }

    public Station getStation() {
        return station;
    }

    public void setStation(Station station) {
        this.station = station;
    }

    public boolean isWorking() {
        return this.getStatusDescription().equalsIgnoreCase("Disponible");
    }

    @Override
    public int hashCode() {
        final int prime = 31;
        int result = 1;
        result = prime * result + ((direction == null) ? 0 : direction.hashCode());
        result = prime * result + ((id == null) ? 0 : id.hashCode());
        result = prime * result + ((situation == null) ? 0 : situation.hashCode());
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
        Elevator other = (Elevator) obj;
        if (direction == null) {
            if (other.direction != null) {
                return false;
            }
        } else if (!direction.equals(other.direction)) {
            return false;
        }
        if (id == null) {
            if (other.id != null) {
                return false;
            }
        } else if (!id.equals(other.id)) {
            return false;
        }
        if (situation == null) {
            if (other.situation != null) {
                return false;
            }
        } else if (!situation.equals(other.situation)) {
            return false;
        }
        return true;
    }

    public int compareTo(Elevator another) {
        return getId().compareTo(another.getId());
    }
}
