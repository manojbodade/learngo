var Dashboard = {};
Dashboard.ServerData = ko.observableArray([]);

function getCategory(val) {
    if(val < 61) return "green";
    else if(val < 81) return "orange";
    else return "red";
}

Dashboard.getData = function() {
    return new Promise(function(res, rej) {
        ajaxPost("/opchadoop/getclusterstatus", {}, function(r){
            res(r.Data);
        }, function() {
            rej("Failed getData");
        })
    });
}
Dashboard.init = function() {
    Dashboard.getData()
             .then(function(data) {
                 var _data = _.map(data, function(d) {
                     var cpuPct = d.CPU.Current.toFixed(2);
                     var storagePct = (d.Storage.Pct * 100).toFixed(2);
                     var memoryPct = 0;
                     var memoryCurr = 0;
                     if(d.Memory.Timeline.length > 0){
                        memoryPct = (d.Memory.Timeline[d.Memory.Timeline.length - 1].Val / d.Memory.Max * 100).toFixed(2);
                        memoryCurr = autoConvertKByte(d.Memory.Timeline[d.Memory.Timeline.length - 1].Val);
                     }
                     var ioPct = 0;
                     return {
                        id: d.ClusterId,
                        name: d.ClusterName,
                        cpu: {
                            pct: cpuPct,
                            amber: d.CPU.Amber,
                            red: d.CPU.Red,
                            green: d.CPU.Green,
                        },
                        memory: {
                            pct: memoryPct,
                            current: memoryCurr,
                            max: autoConvertKByte(d.Memory.Max),
                            amber: d.Memory.Amber,
                            red: d.Memory.Red,
                            green: d.Memory.Green,
                        },
                        storage: {
                            pct: storagePct,
                            used: convertGig(d.Storage.Used, d.Storage.Label)+" "+d.Storage.Label,
                            free: convertGig(d.Storage.Free, d.Storage.Label)+" "+d.Storage.Label,
                            amber: d.Storage.Amber,
                            red: d.Storage.Red,
                            green: d.Storage.Green,
                        },
                        io: {
                            pct: ioPct
                        }
                     }
                 });
                 Dashboard.ServerData(_data);
             })
}

$(function() {
    Dashboard.init();
    $("[data-toggle='tooltip']").tooltip();
})