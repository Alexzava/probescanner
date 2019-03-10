var command = "get_devices";
var ws = new WebSocket("ws://127.0.0.1:8080/scan");

ws.onopen = function() {
	ws.send(command);
};

ws.onmessage = function (evt) { 
	var received_msg = evt.data;
	ParseDevices(received_msg)
};

function ParseDevices(jsonString) {
	// Read filters
	var url = new URL(window.location.href);
	var filterType = url.searchParams.get("f");
	var filterContent = url.searchParams.get("c");

	var devices = JSON.parse(jsonString);
	var table = document.getElementById("devicesTable");
	var rows = 1;
	var totalDevices = 0;

	for(var key in devices) {
		totalDevices++;

	    // Calculate distance (NOT ACCURATE	)
	        // Formula
            // TxPower = -54 (Estimated)
            // N = 5 (2 to 5 depends on location)
            // d = 10^((TxPower - RSSI)/(10*N))
        var TxPower = -54;
        var N = 5;
        var distance = Number(Math.pow(10, ((TxPower - devices[key].RSSI)/(10*N)))).toFixed(1);

		// Check filters
		if(filterType != null) {
			if(filterType == "mac" && devices[key].MAC != filterContent) {
				continue;
			} else if(filterType == "bssid" && devices[key].BSSID != filterContent) {
				continue;
			} else if(filterType == "rssimax" && devices[key].RSSI > filterContent) {
				continue;
			} else if(filterType == "rssimin" && devices[key].RSSI < filterContent) {
				continue;
			} else if(filterType == "vendor" && devices[key].MAC.substring(0, 8) != filterContent) {
				continue;
			} else if(filterType == "distmin" && distance > filterContent) {
				continue;
			} else if(filterType == "distmax" && distance < filterContent) {
				continue;
			}
		}

		// Insert device to table
		var row = table.insertRow(rows);
	    var cellMAC = row.insertCell(0);
	    var cellVendor = row.insertCell(1);
	    var cellRSSI = row.insertCell(2);
	    var cellDistance = row.insertCell(3);
	    var cellBSSID = row.insertCell(4);
	    var cellTime = row.insertCell(5);

	    cellMAC.innerHTML = '<a href="?f=mac&c='+devices[key].MAC+'">'+devices[key].MAC+'</a>'; // devices[key].MAC
	    cellVendor.innerHTML = '<a href="?f=vendor&c='+devices[key].MAC.substring(0,8)+'">'+devices[key].Vendor+'</a>';
	    cellRSSI.innerHTML = devices[key].RSSI;
	    cellBSSID.innerHTML = '<a href="?f=bssid&c='+devices[key].BSSID+'">'+devices[key].BSSID+'</a>';

        cellDistance.innerHTML = distance;

	    // Parse time from Unix
	    var date = new Date(devices[key].DetectedTime*1000)
	    cellTime.innerHTML = date.toString();

	    row++;
	}

	document.getElementById('devicesCounter').innerHTML = totalDevices;
}