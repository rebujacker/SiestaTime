  
//// Index.JS will hold support client side functions

//Refresh every buffer of data in first opening,a nd then keep refreshing on key clicks
  var jobs = "";
  var logs = "";
  var implants = "";
  var redirectors = "";
  var bichitos = "";
  var vps = "";
  var domains = "";
  var stagings = "";
  var reports = "";
  
//AJAX functions to pull data from client server

function getJobs() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/jobs",
        success: function (response){
          if (response == null){
            jobs = "{}";
          }else{
            jobs = response;
          }
        }        
    });
}

function getLogs() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/logs",
        success: function (response){
          if (response == null){
            logs = "";
          }else{
            logs = response;
          }
        }        
    });
}



function getImplants() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/implants",
        success: function (response){
          if (response == null){
            implants = "";
          }else{
            implants = response;
          }
        }        
    });
}


function getVps() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/vps",
        success: function (response){
          if (response == null){
            vps = "";
          }else{
            vps = response;
          }
        }        
    });
}

function getDomains() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/domains",
        success: function (response){
          if (response == null){
            domains = "";
          }else{
            domains = response;
          }
        }        
    });
}

function getStagings() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/stagings",
        success: function (response){
          if (response == null){
            stagings = "";
          }else{
            stagings = response;
          }
        }        
    });
}

function getReports() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/reports",
        success: function (response){
          if (response == null){
            reports = "";
          }else{
            reports = response;
          }
        }        
    });
}

function getRedirectors() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/redirectors",
        success: function (response){
          if (response == null){
            redirectors = "";
          }else{
            redirectors = response;
          }
        }        
    });
}

function getBichitos() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/bichitos",
        success: function (response){
          if (response == null){
            bichitos = "";
          }else{
            bichitos = response;
          }
        }        
    });
}


//Some  other support functions

//Get an array with all Implants Rid/Bid's
function getImpantIds(implantName) {
  var returnArray = [];
  for (var i = 0; i < implants.length; i++){
    var row = implants[i];
    if (row.name == implantName){
      returnArray.push(row.pid)
    }
  }
  
  return returnArray
}

  // Not working on opening STIme research
  getJobs();
  getLogs();
  getImplants();
  getRedirectors();
  getBichitos();
  getVps();
  getDomains();
  getStagings();
  getReports();
  $(".STheader").load('./components/header/header.html');

$(document).ready(function() {
        
  $("#sidebar-menu").on("click","a",function() {
        
        //Prepare Menu elements
        var link = $(this);
        var closestUpper_ul = link.closest("ul");
        //var parallel_active_links = closestUpper_ul.find(".active")
        var closestUpper_li = link.closest("li");
        var lower_ul = closestUpper_li.children("ul")
        //var link_status = closestUpper_li.hasClass("active");
        
        // On Menu click, refresh Memory Data
        getJobs();
        getLogs();
        getImplants();
        getRedirectors();
        getBichitos();
        getVps();
        getDomains();
        getStagings();
        getReports();
        
        $(".STheader").load('./components/header/header.html');

        switch(link.attr("class")) {
          
          //Hive Menu Options
          case "hivJobs":
            $(".STmain").empty();
            $(".STmain").attr("id","Hive");
            $(".STmain").load('./components/jobs/jobs.html')
            break;   
          case "hivLogs":
            $(".STmain").empty();
            $(".STmain").attr("id","Hive");
            $(".STmain").load('./components/logs/logs.html')
            break;   

          //Implant, bots sliding menu
          case "implantList":
            lower_ul.empty();
            for (i = 0; i < implants.length; i++){
              //console.log("Adding implat");
              var row = implants[i];
              lower_ul.append("<li class=\"implantli\"><a href=\"#\" class=\"implant\" id=\""+row.name+"\">"+row.name+"<span class=\"fa fa-chevron-down\"></span></a><ul class=\"nav child_menu\" style=\"display: block;\"></ul></li>");
            }
            break;
          case "implant":
            $(".STmain").empty();
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/implant/implant.html')
            lower_ul.empty();
            //Get Parent Implant RespTime, compare it with "lastchecked" to see if it is online
            var bidentities = [];
            for (i = 0; i < bichitos.length; i++){
              var row = bichitos[i];
              var binfo = row.info;
              //console.log(binfo);
              var infoJson = JSON.parse(row.info)
              //console.log(bidentity);
              //Client side redirector organization
              if (row.implantname == link.attr("id")){
                var bidentity = infoJson.mac.replace(/\n/g, '')+infoJson.hostname.replace(/\n/g, '');
                if (bidentities.includes(bidentity)){
                  continue;
                }
                bidentities.push(bidentity);
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"host\" id=\""+bidentity+"\">"+infoJson.hostname.replace(/\n/g, '')+"<span class=\"fa fa-chevron-down\"></span></a><ul class=\"nav child_menu\" style=\"display: block;\"></ul></li>");
              }
            }
            break;
          case "host":
            //$(".STmain").attr("id",link.attr("id"));
            //$(".STmain").load('./components/host/host.html')
            lower_ul.empty();
            //Get Parent Implant RespTime, compare it with "lastchecked" to see if it is online
            for (i = 0; i < bichitos.length; i++){
              var row = bichitos[i];
              var binfo = row.info;
              //console.log(binfo);
              var infoJson = JSON.parse(row.info)
              var bidentity = infoJson.mac.replace(/\n/g, '') + infoJson.hostname.replace(/\n/g, '');
              //console.log("LInk:"+link.closest('.implantli').find('.implant').html());
              //console.log(link.attr("id"));
              //Client side redirector organization
              if ((row.implantname == link.closest('.implantli').find('.implant').attr("id")) && (bidentity == link.attr("id"))){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"bichito\" id=\""+row.bid+"\">"+row.bid+"</a></li>");
              }
            }
            break;
          case "bichito":
            $(".STmain").attr("id",link.attr("id"));
            //console.log($(".STmain").attr("id",link.attr("id")));
            $(".STmain").empty();
            $(".STmain").load('./components/bichito/bichito.html')
            lower_ul.empty();
            break;

          // Vps,domains and Staging List/Menus
          case "vpsList":
            lower_ul.empty();
            for (i = 0; i < vps.length; i++){
              var row = vps[i];
              lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"vps\" id=\""+row.name+"\">"+row.name+"</a></li>");
            }
            break;
          case "vps":
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/vps/vps.html')
            break;
          case "domainsList":
            lower_ul.empty();
            for (i = 0; i < domains.length; i++){
              var row = domains[i];
              if (row.dtype != "gmail"){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"domain\" id=\""+row.name+"\">"+row.name+"</a></li>");
              }
            }
            break;
          case "saasList":
            lower_ul.empty();
            for (i = 0; i < domains.length; i++){
              var row = domains[i];
              if (row.dtype == "gmail"){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"domain\" id=\""+row.name+"\">"+row.name+"</a></li>");
              }            }
            break;
          case "domain":
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/domain/domain.html')
            break;
          case "dropletList":
            lower_ul.empty();
            for (i = 0; i < stagings.length; i++){
              var row = stagings[i];
              if (row.stype == "https_droplet_letsencrypt"){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"staging\" id=\""+row.name+"\">"+row.name+"</a></li>");
              }
            }
            break;
          case "hanlerList":
            lower_ul.empty();
            for (i = 0; i < stagings.length; i++){
              var row = stagings[i];
              if (row.stype != "https_droplet_letsencrypt"){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"staging\" id=\""+row.name+"\">"+row.name+"</a></li>");
              }
            }
            break;
          case "staging":
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/staging/staging.html')
            break;
          case "reportsList":
            lower_ul.empty();
            for (i = 0; i < reports.length; i++){
              var row = reports[i];
              lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"report\" id=\""+row.name+"\">"+row.name+"</a></li>");
            }
            break;
          case "report":
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/report/report.html')
            break;

          //Creation Menus
          case "createimplant":
            $(".STmain").load('./components/createforms/createImplant.html')
            break;
          case "createvps":
            $(".STmain").load('./components/createforms/createVPS.html')
            break;
          case "createdomain":
            $(".STmain").load('./components/createforms/createDomain.html')
            break;
          case "createstaging":
            $(".STmain").load('./components/createforms/createStaging.html')
            break;
          case "createreport":
            $(".STmain").load('./components/createforms/createReport.html')
            break;
          default:
        }
    })
})









