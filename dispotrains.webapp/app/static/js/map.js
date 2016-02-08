function AccessMap(mapEl, captionEl) {
  this.accessible = false;
  this.statistics = false;
  this._mapEl = mapEl;
  this._captionEl = captionEl;
  this._points = [];
  this._data = undefined;
  this._lastSelectedPoint = undefined;
  this._d3Voronoi = d3.geom.voronoi().x(function(d) {
    return d.x;
  }).y(function(d) { return d.y; });
  this._setupMapbox();
}

AccessMap._DISPOTRAINS_STATIONS =
    "http://dispotrains.membrives.fr/app/GetStations/";
AccessMap._DISPOTRAINS_STATS = "http://dispotrains.membrives.fr/app/AllStats/";
AccessMap._STATIONS_CSV = "full-list.csv";

AccessMap.prototype._setupMapbox = function() {
  L.mapbox.accessToken =
      'pk.eyJ1IjoiZW1lbWJyaXZlcyIsImEiOiIwNDViZWQyODJhNTczNTg4ZWEzNzI4MzllNzk4ODk1NyJ9.ijO7LzQGt_kX1IwAOrUYzA';
  this._map =
      L.mapbox.map(this._mapEl, 'emembrives.d5f86755')
          .fitBounds([ [ 49.241299, 3.55852 ], [ 48.120319, 1.4467 ] ]);

  var self = this;

  var mapLayer = {
    onAdd : function(map) {
      map.on('viewreset moveend', function() { self.loadAndDraw(); });
      self.loadAndDraw();
    }
  };

  this._map.on('ready', function() {
    self.loadAndDraw();
    self._map.addLayer(mapLayer);
  });
};

AccessMap.prototype._getData = function() {
  var availabilityPromise =
      d3.promise.json(AccessMap._DISPOTRAINS_STATIONS)
          .then(function(stations) {
            for (var j = 0; j < stations.length; j++) {
              var d = stations[j];
              var good = true;
              for (var i = 0; i < d.elevators.length; i++) {
                if (d.elevators[i].status.state != "Disponible") {
                  good = false;
                  break;
                }
              }
              d.good = good;
            }
            return stations;
          });
  var statisticsPromise = d3.promise.json(AccessMap._DISPOTRAINS_STATS);
  var stationsPromise = d3.promise.csv(AccessMap._STATIONS_CSV);
  return Promise
      .all([ availabilityPromise, stationsPromise, statisticsPromise ])
      .then(this._mergeData);
};

AccessMap.prototype._mergeData = function(values) {
  var availabilities = {};
  for (var i = 0; i < values[0].length; i++) {
    availabilities[values[0][i].name.toLowerCase()] = values[0][i];
  }
  var stations = values[1];
  var statistics = {};
  for (var i = 0; i < values[2].length; i++) {
    statistics[values[2][i].name.toLowerCase()] = values[2][i];
  }

  var merged_stations = stations.map(function(d, index) {
    d.accessible = d.accessible === "True";
    if (!d.accessible) {
      return d;
    }
    if (d.name.toLowerCase() in availabilities) {
      var key = d.name.toLowerCase();
      d.name = availabilities[key].displayname;
      d.dispotrains_id = availabilities[key].name;
      d.good = availabilities[key].good;
      d.lines = availabilities[key].lines;
      d.elevators = availabilities[key].elevators;
      d.percentfunction = statistics[key].percentfunction;
      return d;
    }
    console.log("Unable to merge station " + d.name);
    console.log(d);
    return d;
  });
  return merged_stations;
};

AccessMap.prototype.loadAndDraw = function() {
  var self = this;
  if (this._data === undefined) {
    this._data = this._getData();
  }
  return this._data.then(function(p) { self.draw(p); });
};

var metersPerPixel = function(latitude, zoomLevel) {
  var earthCircumference = 40075017;
  var latitudeRadians = latitude * (Math.PI / 180);
  return earthCircumference * Math.cos(latitudeRadians) /
         Math.pow(2, zoomLevel + 8);
};

var pixelValue = function(latitude, meters, zoomLevel) {
  return meters / metersPerPixel(latitude, zoomLevel);
};

AccessMap.prototype.draw = function(points) {
  var self = this;
  d3.select('#overlay').remove();

  var bounds = this._map.getBounds();
  var topLeft = this._map.latLngToLayerPoint(bounds.getNorthWest());
  var bottomRight = this._map.latLngToLayerPoint(bounds.getSouthEast());
  var existing = d3.set();
  var drawLimit = bounds.pad(0.4);

  filteredPoints = points.filter(function(d, i) {
    var latlng = new L.LatLng(d.latitude, d.longitude);

    if (!drawLimit.contains(latlng)) {
      return false
    };

    if (!d.accessible && self.accessible) {
      return false;
    }

    var point = self._map.latLngToLayerPoint(latlng);

    key = point.toString();
    if (existing.has(key)) {
      return false
    };
    existing.add(key);

    d.x = point.x;
    d.y = point.y;
    return true;
  });

  var svg = d3.select(this._map.getPanes().overlayPane)
                .append("svg")
                .attr('id', 'overlay')
                .attr("class", "leaflet-zoom-hide")
                .style("width", this._map.getSize().x + 'px')
                .style("height", this._map.getSize().y + 'px')
                .style("margin-left", topLeft.x + "px")
                .style("margin-top", topLeft.y + "px")
                .append("g")
                .attr("transform",
                      "translate(" + (-topLeft.x) + "," + (-topLeft.y) + ")");

  var clips = svg.append("svg:g").attr("id", "point-clips");
  var points = svg.append("svg:g").attr("id", "points");
  var paths = svg.append("svg:g").attr("id", "point-paths");

  clips.selectAll("clipPath")
      .data(filteredPoints)
      .enter()
      .append("svg:clipPath")
      .attr("id", function(d, i) { return "clip-" + i; })
      .append("svg:circle")
      .attr('cx', function(d) { return d.x; })
      .attr('cy', function(d) { return d.y; })
      .attr('r', pixelValue(48.8534100, 20000, this._map.getZoom()));

  var datapointFunc = function(datapoint) {
    return "M" + datapoint.join(",") + "Z";
  };
  var areaPath = paths.selectAll("path")
      .data(this._d3Voronoi(filteredPoints))
      .enter()
      .append("svg:path");
  areaPath.attr("d", datapointFunc)
      .attr("id", function(d, i) { return "path-" + i; })
      .attr("clip-path", function(d, i) { return "url(#clip-" + i + ")"; })
      .style("stroke", d3.rgb(50, 50, 50))
      .on('click', function(d) { self._selectPoint(d3.select(this), d); })
      .classed("selected", function(d) { return this._lastSelectedPoint == d })
      .classed("inaccessible", function(d) { return !d.point.accessible; })
      .classed("malfunction", function(d) { return d.point.accessible && (
            !d.point.good || self.statistics); })
      .classed("ok", function(d) {
        return d.point.accessible && d.point.good && !self.statistics; });
  if (self.statistics) {
    areaPath.style("opacity", function(d) {
      return (100 - d.point.percentfunction)/100;
    });
  }

  paths.selectAll("path")
      .on("mouseover",
          function(d, i) { d3.select(this).classed("mouseover", true); })
      .on("mouseout", function(d, i) {
        d3.select(this).classed("mouseover", false);
      });

  points.selectAll("circle")
      .data(filteredPoints)
      .enter()
      .append("svg:circle")
      .attr("id", function(d, i) { return "point-" + i; })
      .attr("transform",
            function(d) { return "translate(" + d.x + "," + d.y + ")"; })
      .attr("r", 1.5)
      .attr('stroke', 'none');
};

AccessMap.prototype._selectPoint = function(cell, point) {
  d3.selectAll('.selected').classed('selected', false);

  if (this._lastSelectedPoint == point) {
    this._lastSelectedPoint = null;
    d3.select('#selected').classed('hidden', true);
    return;
  }

  this._lastSelectedPoint = point;
  cell.classed('selected', true);

  d3.select('#selected').classed('hidden', false);
  d3.select('#selected #header').text(point.point.name);

  if (!point.point.accessible) {
    d3.select("#selected #card-inaccessible").style("display", null);
    d3.select("#selected #card-accessible").style("display", "none");
    return;
  } else {
    d3.select("#selected #card-inaccessible").style("display", "none");
    d3.select("#selected #card-accessible").style("display", null);
  }
  d3.select("#selected #line-count").text(point.point.lines.length);
  d3.select("#selected #lines").text(this._lineStr(point.point.lines));
  d3.select("#selected #elevator-count")
      .text(point.point.elevators.length)
          d3.select("#selected #broken-elevator-count")
      .text(this._elevatorStr(point.point.elevators));
  d3.select("#selected #function").text(point.point.percentfunction.toFixed(1));
  d3.select("#selected a#dispotrains")
      .attr("href", "/gare/" + point.point.dispotrains_id);
};

AccessMap.prototype._lineStr = function(lines) {
  var s = lines[0].id;
  for (var i = 1; i < lines.length; i++) {
    s += ", ";
    s += lines[i].id;
  }
  return s;
};

AccessMap.prototype._elevatorStr = function(elevators) {
  var bad = 0;
  for (var i in elevators) {
    if (elevators[i].status.state !== "Disponible") {
      bad++;
    }
  }
  return bad;
};
