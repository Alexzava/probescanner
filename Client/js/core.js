var command = "get_devices";
var ws = new WebSocket("ws://127.0.0.1:8080/scan");

setInterval(function() {
	ws.send(command);
}, 5000);


/*ws.onopen = function() {
	ws.send(command);
};*/

ws.onmessage = function (evt) { 
	var received_msg = evt.data;
	ParseDevices(received_msg)
};

function ParseDevices(jsonString) {
	var devices = JSON.parse(jsonString);
	var table = document.getElementById("devicesTable");
	var rows = 1;
	var totalDevices = 0;

	for (var key in devices) {
		// Insert device to table
		var row = table.insertRow(rows);
	    var cellMAC = row.insertCell(0);
	    var cellVendor = row.insertCell(1);
	    var cellRSSI = row.insertCell(2);
	    var cellBSSID = row.insertCell(3);
	    var cellTime = row.insertCell(4);

	    cellMAC.innerHTML = devices[key].MAC;
	    cellVendor.innerHTML = devices[key].Vendor;
	    cellRSSI.innerHTML = devices[key].RSSI;
	    cellBSSID.innerHTML = devices[key].BSSID;

	    // Parse time from Unix
	    var date = new Date(devices[key].DetectedTime*1000)
	    cellTime.innerHTML = date.toString();

	    row++;
	    totalDevices++;
	}

	document.getElementById('devicesCounter').innerHTML = totalDevices;
}