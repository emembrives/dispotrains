function getDateDisplay(dateObj) {
     var date = moment(dateObj);
     return date.format("llll");
}

    function displayStation(urlObj, options) {
        var stationId = urlObj.hash.replace(/.*id=/, "");
        var main = $('#station');
        var header = main.children( ":jqmData(role=header)" );
        var content = main.children( ":jqmData(role=content)" );
        $.mobile.showPageLoadingMsg();

        $.getJSON('/app/GetStation/' + stationId + '/', function(data) {
            header.find( "h1" ).html(data.displayname);

            var linesDiv = main.find('#lines').find(".data").html(''),
                elevatorsDiv = main.find('#elevators').find(".data").html('');

            var linesUl = $('<ul>').appendTo(linesDiv).attr("data-role", "listview").attr("data-inset", "false").attr("data-filter", "false");
            $.each(data["lines"], function(index, line) {
                var li = $("<li></li>");
                var a = $("<a></a>").appendTo(li);
                a.id = line["id"];
                a.attr("href", "#ligne?id=" + line["id"]);
                a.append(line["network"] + " " + line["id"]);
                li.appendTo(linesUl);
            });

            elevatorsUl = $('<ul>').appendTo(elevatorsDiv).attr("data-role", "listview").attr("data-inset", "false").attr("data-filter", "false");
            $.each(data["elevators"], function(index, elevator) {
                var li = $("<li></li>");
                li.append('<h3 style="white-space: normal">' + elevator["situation"] + "</h3>");

                if (elevator["code"] != "rampe") {
                    li.append("<p>" + elevator["direction"] + ' - Situation en date du <strong>' + getDateDisplay(elevator.status.lastupdate) + '</strong></p><p class="ui-li-aside">' + elevator.status.state +'</p>');
                } else {
                    li.append("<p>" + elevator["direction"] + "</p>");
                }
                if (elevator.status.state != "Disponible" && elevator["code"] != "rampe") {
                    li.attr("data-icon", "alert");
                    li.attr("data-theme", "e");
                }
                li.appendTo(elevatorsUl);
            });

            main.page();
            content.find( ":jqmData(role=listview)" ).listview();

            // We don't want the data-url of the page we just modified
            // to be the url that shows up in the browser's location field,
	    	// so set the dataUrl option to the URL for the category
		    // we just loaded.
		    options.dataUrl = urlObj.href;

            // Do the page change
            $.mobile.hidePageLoadingMsg();
            $.mobile.changePage( main, options );
        });
    }

    function displayStationsByLine(urlObj, options) {
        var lineId = urlObj.hash.replace(/.*id=/, "");

        var app = $("#app");
        var main = $('#line');

        var header = main.children( ":jqmData(role=header)" );
        var content = main.children( ":jqmData(role=content)" );

        $.mobile.showPageLoadingMsg();

        loadLinesData();
        var getLines = app.data('getLines');
        getLines.success(function(lines) {
            content.html('');

            for (var i = 0; i < lines.length; i++) {
                if (lines[i].id == lineId) {
                    var data = lines[i];
                    break;
                }
            }

            var network = data["network"];

            header.find( "h1" ).html(network + " " + lineId);

            var ul = $('<ul>').appendTo(content).attr("data-role", "listview").attr("data-inset", "true").attr("data-filter", "true");

            $('<li></li>').attr("data-role", "list-divider").html("Avec dysfonctionnement(s)").appendTo(ul);
            $.each(data.badStations, function(index, station) {
                var li = $("<li></li>");
                var a = $("<a></a>").appendTo(li);
                a.id = station["name"];
                a.attr("href", "#station?id=" + station["name"]);
                li.attr("data-icon", "alert");
                li.attr("data-theme", "e");
                a.append(station["displayname"]);
                li.appendTo(ul);
            });

            $('<li></li>').attr("data-role", "list-divider").html("En fonctionnement").appendTo(ul);
            $.each(data["goodStations"], function(index, station) {
                var li = $("<li></li>");
                var a = $("<a></a>").appendTo(li);
                a.id = station["name"];
                a.attr("href", "#station?id=" + station["name"]);
                a.append(station["displayname"]);
                li.appendTo(ul);
            });
            main.data("lineId", lineId);
            main.data("update", new Date());
            main.page();
            content.find( ":jqmData(role=listview)" ).listview();

            // We don't want the data-url of the page we just modified
    		// to be the url that shows up in the browser's location field,
	    	// so set the dataUrl option to the URL for the category
		    // we just loaded.
		    options.dataUrl = urlObj.href;

            // Do the page change
            $.mobile.hidePageLoadingMsg();
            $.mobile.changePage( main, options );
        });
    }

    function loadLinesData() {
        var app = $("#app");
        var getLines = $.getJSON('/app/GetLines/', function(data) {});
        app.data('getLines', getLines);
    }

    function displayLines(urlObj, options) {
        var app = $("#app");
        var main = $("#main");
        var content = main.children( ":jqmData(role=content)" );
        if (content.html() != "") {
            // Do the page change
            $.mobile.changePage( main, options );
            return;
        }

        loadLinesData();
        var getLines = app.data('getLines');
        getLines.success(function(data) {
            content.html('');
            var ul = $('<ul>').appendTo(content).attr("data-role", "listview").attr("data-inset", "true").attr("data-filter", "false");

            $.each(data, function(index, ligne) {
                var li = $("<li></li>");
                var a = $("<a></a>").appendTo(li);
                a.id = ligne["id"];
                a.attr("href", "#line?id=" + ligne["id"]);
                a.append(ligne["network"] + ' ' + ligne["id"]);
                li.appendTo(ul);
            });

            main.page();

            content.find( ":jqmData(role=listview)" ).listview();

        	// We don't want the data-url of the page we just modified
    		// to be the url that shows up in the browser's location field,
	    	// so set the dataUrl option to the URL for the category
		    // we just loaded.
		    options.dataUrl = urlObj.href;

            // Do the page change
            $.mobile.changePage( main, options );
        });
    }

// Listen for any attempts to call changePage().
$(document).bind( "pagebeforechange", function( e, data ) {
    // We only want to handle changePage() calls where the caller is
	// asking us to load a page by URL.
	if ( typeof data.toPage === "string" ) {
		// We are being asked to load a page by URL, but we only
		// want to handle URLs that request the data for a specific
		// category.
		var u = $.mobile.path.parseUrl( data.toPage );
        _gaq.push(['_trackPageview', data.toPage]);
		if ( u.hash.indexOf("#line") == 0 ) {
			// Display one train line
            displayStationsByLine( u, data.options );

			// Make sure to tell changePage() we've handled this call so it doesn't
			// have to do anything.
        	e.preventDefault();
		} else if (u.hash.indexOf("#station") == 0) {
            displayStation(u, data.options);

            e.preventDefault();

		} else if (u.hash == "") {
    	    displayLines(u, data.options);
            // Make sure to tell changePage() we've handled this call so it doesn't
			// have to do anything.
            e.preventDefault();
		}
	}
});

$(document).ready(function() {
    moment.lang('fr');
    displayLines();
});

