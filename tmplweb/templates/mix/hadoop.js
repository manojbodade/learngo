var OPHadoop = {};
OPHadoop.OpcData = ko.observableArray([]);
OPHadoop.Modal = {
    title: ko.observable(""),
    selector: "#modalDetail",
    additional : {
        memory: function(additional) {
            var series = additional.series;
            console.log(parseFloat(series[series.length - 1].Val), parseFloat(additional.max))
            var pct = parseFloat(series[series.length - 1].Val) / parseFloat(additional.max) * 100;
            var html = '<div class="opProgress text-center"><div class="bar default" style="width: '+pct+'%"></div></div>';
            return html;
        }
    }
}

function DetailGrid(id, name, type) {
    var columnGrid = {
        storage: [
            { title: "Nodes", field: "HostName" },
            { title: "Total Storage (GB)", field: "DiskTotal", attributes:{ class: "pdl-20" }, template: function(o) { return o.DiskTotal.toFixed(2) } },
            { title: "Usage Storage (GB)", field: "DiskUsage", attributes:{ class: "pdl-20" }, template: function(o) { return o.DiskUsage.toFixed(2) } },
            { title: "Available Storage (GB)", field: "DiskFree", attributes:{ class: "pdl-20" }, template: function(o) { return o.DiskFree.toFixed(2) } },
            { title: "%", field: "DiskPct", attributes:{ class: "pdl-20" }, template: function(o) { return o.DiskPct.toFixed(2) } },
        ],
        memory: [
            { title: "Nodes", field: "HostName" },
            { title: "Total Memory (GB)", field: "MemTotal", attributes:{ class: "pdl-20" }, template: function(o) {return kb2gb(o.MemTotal)} },
            { title: "Used (GB)", field: "MemUsed", attributes:{ class: "pdl-20" }, template: function(o) {return kb2gb(o.MemUsed)} },
            { title: "Available (GB)", field: "MemFree", attributes:{ class: "pdl-20" }, template: function(o) {return kb2gb(o.MemFree)} },
            { title: "%", field: "MemPct", attributes:{ class: "pdl-20" }, template: function(o) {return o.MemPct.toFixed(2)} },
            {
                title: "Aging of Running Jobs",
                columns: [
                    { title: "> 60 mins", field: "" },
                    { title: "30 - 60 mins", field: "" },
                    { title: "15 - 30 mins", field: "" },
                    { title: "< 15 mins", field: "" },
                ]
            }
        ]
    }
    var self = this;
    self.selector = "#gridModal";
    self.render = function() {
        return new Promise(function(res, rej) {
            $(self.selector).html("");
            $(self.selector).kendoGrid({
                dataSource: {
                    transport: {
                        read: function(opt) {
                            ajaxPost("/nodeinfo/statusinfo", { ClusterId: id }, function(r) {
                                opt.success({ data: r.Data, total: r.Data.length });
                                res("success");
                            }, function() {
                                rej("failed get data");
                            })
                        }
                    },
                    schema: {
                        data: "data",
                        total: "total"
                    },
                    pageSize: 10,
                },
                pageable: true,
                columns: columnGrid[type]
            })
        });
    }
}

function OpcBarChart(name, id, data) {
    var self = this;
    self.pct = ko.observable(0);
    self.diskUsedInf = ko.observable("");
    self.diskFree = ko.observable("");
    self.type = ko.observable("");

    self.render = function() {
        var pct = data.Pct.toFixed(2);
        var free = data.Free.toFixed(2);
        var used = data.Used.toFixed(2);
        if(pct > 85) self.type("bar red");
        else if(pct > 70) self.type("bar yellow");
        else self.type("bar green");
        self.pct(pct);
        self.diskUsedInf(used+data.Label+" ("+pct+"%)");
        self.diskFree(+free+data.Label+" ("+(100 - pct)+"%)");
    }
}

function OpcAreaChart(name, id, data) {
    var self = this;
    self.process = ko.observable("");
    self.selector = "#opChartArea-"+id

    self.formatData = function(data) {
        var newFormat = [];
        var colors = ["rgb(0, 117, 176)", "rgb(0, 92, 132)", "rgb(63, 156, 53)", "rgb(105, 190, 40)"];
        var keys = Object.keys(data);
        for(var i = 0; i < keys.length; i++) {
            newFormat.push({
                name: keys[i],
                data: data[keys[i]],
                color: colors[i]
            });
        }
        return newFormat;
    }
    self.render = function() {
        self.process("Current: "+data.Process+" processes");

        var series = self.formatData(data.Data);
        var selector = $(self.selector).find(".area");
        selector.html("");
        selector.kendoSparkline({
            seriesDefaults: {
                type: "area",
                field: "Val"
            },
            chartArea: {
                height: 125
            },
            series: series,
            tooltip: {
                template: function(o) {
                    var data = o.dataItem;
                    var d = new Date(0);
                    d.setUTCSeconds(data.Ts);
                    return data.Val+"% at "+moment(d).format("HH:mm:ss");
                }
            },
        })
    }
}

function OpcLineChart(name, id ,data) {
    var self = this;
    self.maxInf = ko.observable("");
    self.max = 0;
    self.selector = "#opChartLine-"+id;
    self.series = null;

    self.formatData = function(data) {
        return _.map(data, function(d) {
            return {
                Val: (d.Val / 1048576).toFixed(2),
                Ts: d.Ts,
                Tst: d.Tst
            }
        })
    }
    self.formatCategories = function(data) {
        return _.map(data, function(d) {
            var date = new Date(0);
            date.setUTCSeconds(d.Ts);
            return moment(date).format("HH:mm");
        })
    }

    self.render = function() {
        var max = (data.Max / 1048576).toFixed(2);
        self.max = max;
        self.maxInf("Max: "+max+"GB");
        var selector = $(self.selector).find(".line");
        var categories= self.formatCategories(data.Timeline);
        selector.html("");
        self.series = self.formatData(data.Timeline);
        var dataPeak = _.maxBy(self.series, function(o) { return o.Val });
        selector.kendoSparkline({
            data: self.series,
            plotAreaClick: function() {
                OPHadoop.handleChartClick(id, name, "memory", { series: self.series});
            },
            seriesDefaults: {
                type: "line",
                line: {
                    width: 2
                },
                field: "Val",
                color: "rgb(40, 144, 192)"
            },
            chartArea: {
                height: 125,
            },
            valueAxis:{
                visible: true,
                majorUnit: Math.round((parseFloat(dataPeak.Val) / 3) * 100) / 100,
                plotBands: [
                    {
                        from: max,
                        to: max + 2,
                        color: "#000",
                    }
                ],
                labels:{
                    visible: true,
                    font: "10px Arial,Helvetica,sans-serif"
                }
            },
            categoryAxis: {
                visible: true,
                categories: categories,
                majorTicks: {
                    step: categories.length - 1,
                },
                line: {
                    visible: true,
                },
                labels: {
                    step: categories.length - 1,
                    font: "10px Arial,Helvetica,sans-serif"
                }
            },
            tooltip: {
                template: function(o) {
                    var data = o.dataItem;
                    var d = new Date(0);
                    d.setUTCSeconds(data.Ts);
                    return data.Val+"GB at "+moment(d).format("HH:mm:ss");
                }
            },
        })
    }
}

function OpcCircleChart(name, id, data) {
    var self = this;
    self.selector = "#opChartCircle-"+id;

    self.getCircleChart = function(pct) {
        return '<svg class="circle-chart" viewbox="0 0 33.83098862 33.83098862" width="150" height="150" xmlns="http://www.w3.org/2000/svg"><circle class="circle-chart__background" stroke="rgb(195, 226, 193)" stroke-width="2" fill="none" cx="16.91549431" cy="16.91549431" r="15.91549431" /><circle class="circle-chart__circle" stroke="rgb(63, 156, 53)" stroke-width="2" stroke-dasharray="'+pct+',100" stroke-linecap="round" fill="none" cx="16.91549431" cy="16.91549431" r="15.91549431" /><g class="circle-chart__info"><text class="circle-chart__percent" x="16.91549431" y="16" alignment-baseline="central" text-anchor="middle" fill="rgb(106, 193, 123)" font-size="6">'+pct+'%</text></g></svg>';
    }
    self.render = function() {
        var html = self.getCircleChart(data.Current.toFixed(2));
        $(self.selector).html(html);
    }
}

OPHadoop.getData = function() {
    return new Promise(function(res, rej) {
        ajaxPost("/opchadoop/getclusterstatus",{}, function(r) {
            res(r.Data);
        }, function() {
            rej("Failed get data");
        })
    });
}
OPHadoop.buildData = function(data) {
    var newFormat = [];
    for(var i = 0; i < data.length; i++) {
        var d = data[i];
        newFormat.push({
            name: d.ClusterName,
            id: d.ClusterId,
            storage: new OpcBarChart(d.ClusterName, d.ClusterId, d.Storage),
            io: new OpcAreaChart(d.ClusterName, d.ClusterId, d.IOinfo),
            memory: new OpcLineChart(d.ClusterName, d.ClusterId, d.Memory),
            cpu: new OpcCircleChart(d.ClusterName, d.ClusterId, d.CPU),
        })
    }
    return newFormat;
}

OPHadoop.handleChartClick = function(id, name, type, additional) {
    var grid = new DetailGrid(id, name, type);
    grid.render()
        .then(function() {
            OPHadoop.Modal.title(capitalize(type) +" | "+ name);
            $(OPHadoop.Modal.selector).find(".additional").html("");
            if(additional) {
                $(OPHadoop.Modal.selector).find(".additional").html(OPHadoop.Modal.additional[type](additional));
            }
            $(OPHadoop.Modal.selector).modal("show");
        })
        .catch(function(msg) {
            alert(msg);
        });
}

OPHadoop.render = function(interval) {
    OPHadoop.getData()
    .then(function(data) {
        var newData = OPHadoop.buildData(data);
        OPHadoop.OpcData(newData);
    })
    .then(function() {
        var data = OPHadoop.OpcData();
        _.each(data, function(d) {
            d.storage.render();
            d.io.render();
            d.memory.render();
            d.cpu.render();
        })
    })
    .catch(function(msg) {
        if(interval) clearInterval(interval);
        alert(msg);
    })
}
OPHadoop.init = function() {
    OPHadoop.render();
    // var interval = setInterval(function() {
    //     OPHadoop.render(interval);
    // }, 60000);
}

$(function() {
    OPHadoop.init();
})
