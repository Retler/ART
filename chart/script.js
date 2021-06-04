var chartColors = {
    red: 'rgb(255, 99, 132)',
    orange: 'rgb(255, 159, 64)',
    yellow: 'rgb(255, 205, 86)',
    green: 'rgb(75, 192, 192)',
    blue: 'rgb(54, 162, 235)',
    purple: 'rgb(153, 102, 255)',
    grey: 'rgb(201, 203, 207)'
};

window.onload = function() {
    var ctx = document.getElementById('myChart').getContext('2d');
    window.myChart = new Chart(ctx, config);
    addData(window.myChart, Date.now(), 10);
    setTimeout(function(){ 
	addData(window.myChart, Date.now(), 30);
    }, 3000);
};

function addData(chart, label, data) {
    chart.data.labels.push(label);
    chart.data.datasets.forEach((dataset) => {
        dataset.data.push(data);
    });
    chart.update();
}

var color = Chart.helpers.color;

var config = {
    type: 'line',
    data: {
	datasets: [{
	    data: [],
	    borderDash: [8, 4],
	    lineTension: 0,
	    fill: false,
	    borderColor: chartColors.red,
	    backgroundColor: color(chartColors.red).alpha(0.5).rgbString(),
	    label: "Dataset 1"
	}]
    },
    options: {
	title: {
	    display: true,
	    text: 'Line chart (hotizontal scroll) sample'
	},
	scales: {
	    xAxes: [{
		type: 'time',
		time: {
		    unit: "second"
		}
	    }],
	    yAxes: [{
		scaleLabel: {
		    display: true,
		    labelString: "value"
		},
		ticks: {
		    suggestedMin: -100,
		    suggestedMax: 100
		}
	    }]
	}
    }
}


function onRefresh(chart) {
	var now = Date.now();
	chart.data.datasets.forEach(function(dataset) {
		dataset.data.push({
			x: now,
			y: randomScalingFactor()
		});
	});
}

function randomScalingFactor() {
	return (Math.random() > 0.5 ? 1.0 : -1.0) * Math.round(Math.random() * 100);
}
