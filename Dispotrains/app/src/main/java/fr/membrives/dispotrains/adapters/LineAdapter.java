package fr.membrives.dispotrains.adapters;

import java.util.List;

import android.content.Context;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.TextView;
import fr.membrives.dispotrains.data.Line;

public class LineAdapter extends ArrayAdapter<Line> {
    public LineAdapter(Context context, List<Line> objects) {
        super(context, 0, objects);
    }

    @Override
    public View getView(int position, View convertView, ViewGroup parent) {
        // Get the data item for this position
        Line line = getItem(position);
        // Check if an existing view is being reused, otherwise inflate the view
        if (convertView == null) {
            convertView = LayoutInflater.from(getContext()).inflate(
                    android.R.layout.simple_list_item_1, parent, false);
        }
        // Lookup view for data population
        TextView lineName = (TextView) convertView.findViewById(android.R.id.text1);
        // Populate the data into the template view using the data object
        lineName.setText(line.getNetwork() + " " + line.getId());
        // Return the completed view to render on screen
        return convertView;
    }

}
