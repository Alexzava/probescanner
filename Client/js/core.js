/*

	Code written in just an hour, it might be better.
*/

setInterval(function(){ UpdateDevices() }, 3000);

function UpdateDevices() {
	let xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
	        ParseDevices(this.responseText)
	    }
	};
	xmlhttp.open("GET", "http://127.0.0.1:8683", true);
	xmlhttp.send();
}

function ParseDevices(jsonString) {
	// Read filters
	let url = new URL(window.location.href);
	let filterType = url.searchParams.get("f");
	let filterContent = url.searchParams.get("c");

	let devices = JSON.parse(jsonString);
	let table = document.getElementById("devicesTable");
	let rows = 1;
	let totalDevices = 0;

	for(let key in devices) {
		totalDevices++;

	    // Calculate distance (NOT ACCURATE!)
	        // Formula
            // TxPower = -54 (Estimated)
            // N = 2 (2 to 5 depends on location)
            // d = 10^((TxPower - RSSI)/(10*N))
        let TxPower = -54;
        let N = 2;
        let distance = Number(Math.pow(10, ((TxPower - devices[key].RSSI)/(10*N)))).toFixed(1);

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

		// Add device to table
		let row = table.insertRow(rows);
	    let cellMAC = row.insertCell(0);
	    let cellVendor = row.insertCell(1);
	    let cellRSSI = row.insertCell(2);
	    let cellDistance = row.insertCell(3);
	    let cellBSSID = row.insertCell(4);
	    let cellTime = row.insertCell(5);

	    cellMAC.innerHTML = '<a href="?f=mac&c='+devices[key].MAC+'">'+devices[key].MAC+'</a>';
	    cellVendor.innerHTML = '<a href="?f=vendor&c='+devices[key].MAC.substring(0,8)+'">'+devices[key].Vendor+'</a>';
	    cellRSSI.innerHTML = devices[key].RSSI;
	    cellBSSID.innerHTML = '<a href="?f=bssid&c='+devices[key].BSSID+'">'+devices[key].BSSID+'</a>';

        cellDistance.innerHTML = distance + " m";

	    // Parse time
	    let date = new Date(devices[key].DetectedTime*1000)
	    cellTime.innerHTML = date.toString();

	    row++;
	}

	document.getElementById('devicesCounter').innerHTML = totalDevices;
}