var Traffic = {};
Traffic.mode = ko.observable("normal"); // normal || expanded
Traffic.ServiceStatus = {
    INIT : "danger",
    INSTALLING : "danger",
    INSTALL_FAILED : "danger",
    INSTALLED : "danger",
    STARTING : "danger",
    STARTED : "good",
    STOPPING : "danger",
    UNINSTALLING : "danger",
    UNINSTALLED : "danger",
    WIPING_OUT : "danger",
    UPGRADING : "danger",
    DISABLE : "danger",
    UNKNOWN : "warning",
}

function GridServer(serverType, selector) {
    var self = this;
    self.type = serverType;
    self.title = serverType == "prod" ? "Production" : "Non Production";
    self.object = [];

    self.data = [];
    self.column = [{
        title: self.title,
        field: "ServerType.value",
        width: 150,
        headerAttributes: { class: "grid-type" },
        headerTemplate : "<a onclick='expand(\""+self.type+"\")'>"+self.title+"</a>",
    }];

    self.getData = function(payload) {
        var prom = new Promise(function(resolve, reject) {
            ajaxPost("/traffic/gettrafficinfo", payload, function(res) {
                self.buildData(res.Data[self.type]);
                self.buildColumns(res.Data[self.type].Header);
                resolve(true);
            }, function() {
                reject(false);
            });
        });
        return prom;
    };
    self.buildData = function(data) {
        var _data = [];
        var header = data.Header != null ? data.Header : [];
        var data = data.Data != null ? data.Data : [];
        for(var i = 0; i < data.length; i++) {
            var row = { ServerType: { type: "", value: data[i].NodeIp } };
            for(var x = 0; x < data[i].Values.length; x++) {
                row[header[x].ComponentName] = {
                    type: data[i].Values[x] == "" ? "" : Traffic.ServiceStatus[data[i].Values[x]],
                    value: data[i].Values[x] == "" ? "" : data[i].NodeIp
                };
            }
            _data.push(row);
        }
        self.data = self.data.concat(_data);
        
    };
    self.buildColumns = function(header) {
        var _columns = [];
        var header = header != null ? header : [];
        for(var i = 0; i < header.length; i++) {
            _columns.push(self.columnTemplate(header[i]));
        }
        self.column = self.column.concat(_columns);
    }
    self.columnTemplate = function(col) {
        return {
            title: col.DisplayName,
            field: col.ComponentName,
            width: 150,
            template: function(t) {
                return "<a onclick='headerClick(\""+t[col.ComponentName].value+"\")' class='server-cell "+t[col.ComponentName].type+"'>"+t[col.ComponentName].value+"</a>";
            },
            headerTemplate : "<a onclick='headerClick(\""+col.ComponentName+"\")'>"+col.DisplayName+"</a>",
            headerAttributes : {
                class: "server-grid"
            }
        }
    };
    self.buildGrid = function() {
        $(selector).kendoGrid({
            dataSource: {
                data: self.data,
                pageSize: 10,
            },
            scrollable: true,
            pageable: true,
            columns: self.column
        });
        return $(selector).data("kendoGrid");
    }
    self.render = function() {
        self.getData({ ClusterType: self.type })
            .then(function() {
                return self.buildGrid();
            })
            .catch(function() {
                alert("failed getData");
            });
    }
}

Traffic.Grid = {
    prod: new GridServer("prod", "#gridPrd"),
    nonprod: new GridServer("nonprod", "#gridNonPrd"),
}


var modalGridData = [
    {
        Service: "HDFS",
        Version: "0.04",
        Vendor: "ATOS",
        LastUpdate: "06/06/2017",
        Status: "Amber",
        Upgrade: "TBD based on version stability"
    },
    {
        Service: "HDFS",
        Version: "0.04",
        Vendor: "ATOS",
        LastUpdate: "06/06/2017",
        Status: "Amber",
        Upgrade: "TBD based on version stability"
    },
    {
        Service: "HDFS",
        Version: "0.04",
        Vendor: "ATOS",
        LastUpdate: "06/06/2017",
        Status: "Amber",
        Upgrade: "TBD based on version stability"
    }
]

function headerClick(type) {
    Traffic.initDetailGrid();
    $("#modalDetail").modal('show');
}

function expand(type) {
    $("."+type).css("display", "block");
    if(type == "prod") {
        $(".nonprod").css("display", "none");
        $(".prod").removeClass("col-md-6");
        $(".prod").addClass("col-md-12");
    }
    if(type == "nonprod") {
        $(".prod").css("display", "none");
        $(".nonprod").removeClass("col-md-6");
        $(".nonprod").addClass("col-md-12");
    }

    var header = $("."+type).find(".grid-type");
    header.addClass("dropdown-header");
    header.html("<input class='dropdownType' />");
    var gridType = [
        { type: "prod", title: "Production" },
        { type: "nonprod", title: "Non Production" },
    ]
    $("."+type).find(".dropdownType").kendoDropDownList({
        dataSource: {
            data: gridType
        },
        value: type,
        dataTextField: "title",
        dataValueField: "type",
        change: function() {
            expand(this._old);
        }
    })
}

Traffic.initDetailGrid = function() {
    $("#gridModal").kendoGrid({
        dataSource: {
            data: modalGridData
        },
        columns: [
            { title: "Service", field: "Service" },
            { title: "Version", field: "Version" },
            { title: "Vendor", field: "Vendor" },
            { title: "Last Update", field: "LastUpdate" },
            { title: "Status", field: "Status" },
            { title: "Upgrade Plan", field: "Upgrade" },
        ]
    })
}

Traffic.generateProductionGrid = function() {
    Traffic.Grid.prod.render();
}
Traffic.generateNonProductionGrid = function() {
    Traffic.Grid.nonprod.render();
}

Traffic.init = function() {
    this.generateNonProductionGrid();
    this.generateProductionGrid();
}

$(function() {
    Traffic.init();
})