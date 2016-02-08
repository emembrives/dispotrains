package fr.membrives.dispotrains.adapters;

import java.util.List;

import android.content.Context;
import android.graphics.Color;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;
import fr.membrives.dispotrains.R;
import fr.membrives.dispotrains.data.Station;

public class StationAdapter extends ArrayAdapter<Station> {
    public StationAdapter(Context context, List<Station> objects) {
        super(context, 0, objects);
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        // Get the data item for this position
        Station station = getItem(position);
        // Check if an existing view is being reused, otherwise inflate the view
        if (convertView == null) {
            convertView = LayoutInflater.from(getContext()).inflate(
                    android.R.layout.simple_list_item_1, parent, false);
        }
        // Lookup view for data population
        TextView stationName = (TextView) convertView.findViewById(android.R.id.text1);
        // Populate the data into the template view using the data object
        stationName.setText(station.getDisplay());
        if (!station.getWorking()) {
            stationName.setBackgroundColor(getContext().getResources().getColor(R.color.problem));
        } else {
            stationName.setBackgroundColor(Color.TRANSPARENT);
        }
        // Return the completed view to render on screen
        return convertView;
    }
}
