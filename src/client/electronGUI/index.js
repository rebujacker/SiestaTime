//// Index.JS will hold support client side functions

//Encode every received input to avoid Script Injections
var htmlencode = require('htmlencode');

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
  var username = "";
  
//AJAX functions to pull data from client server, the data will come back in the shape of Array of JSON Encoded items
//Struncture can be found in: SiestaTime/src/client/clientHivComs.go
//Handlers of these Functions: SiestaTime/src/client/clientGUI.go 

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


function getUsername() {

    $.ajax({
        type: "GET",
        url: "http://127.0.0.1:8000/username",
        success: function (response){
          if (response == null){
            username = "";
          }else{
            username = response;
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



$(document).ready(function() {

  //Load every object on Initialization
  getUsername();
  getImplants();
  getBichitos();
  getRedirectors();
  getBichitos();
  getVps();
  getDomains();
  getStagings();
  getReports();
  

  //Load Header

  $(".STheader").load('./components/header/header.html');
     
  //Listener for clicks in each option of the main Menu   
  $("#sidebar-menu").on("click","a",function() {

        
        //Prepare Menu elements
        var link = $(this);
        var closestUpper_ul = link.closest("ul");
        var closestUpper_li = link.closest("li");
        var lower_ul = closestUpper_li.children("ul")

        
        $(".STheader").load('./components/header/header.html');

        switch(link.attr("class")) {
          
          //Hive Menu Options, just Load components for Hive Jobs and Logs on click
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

          //Implants: Once Clicked, inject new rows per number of implants, using its Name
          case "implantList":
            getImplants();
            getBichitos();
            getRedirectors();

            lower_ul.empty();
            for (i = 0; i < implants.length; i++){
              var row = implants[i];
              lower_ul.append("<li class=\"implantli\"><a href=\"#\" class=\"implant\" id=\""+htmlencode.htmlEncode(row.name)+"\">"+htmlencode.htmlEncode(row.name)+"<span class=\"fa fa-chevron-down\"></span></a><ul class=\"nav child_menu\" style=\"display: block;\"></ul></li>");
            }
            
            break;

          //Implant: Render Implant information on the Right Panel. Inject new rows per "Host", the bot will be identified
          //using Bichitos information (combination of mac Address and hostname) 
          case "implant":
            getBichitos();
            $(".STmain").empty();
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/implant/implant.html')
            lower_ul.empty();
            //Get Parent Implant RespTime, compare it with "lastchecked" to see if it is online
            var bidentities = [];
            for (i = 0; i < bichitos.length; i++){
              var row = bichitos[i];

              try{
              var infoJson = JSON.parse(row.info);
              } catch (e){
                console.log(e);
              }
              //console.log(bidentity);
              //Client side redirector organization
              if (row.implantname == link.attr("id")){
                var bidentity = infoJson.mac.replace(/\n/g, '')+infoJson.hostname.replace(/\n/g, '');
                if (bidentities.includes(bidentity)){
                  continue;
                }
                bidentities.push(bidentity);
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"host\" id=\""+htmlencode.htmlEncode(bidentity)+"\">"+htmlencode.htmlEncode(infoJson.hostname.replace(/\n/g, ''))+"<span class=\"fa fa-chevron-down\"></span></a><ul class=\"nav child_menu\" style=\"display: block;\"></ul></li>");
              }
            }
            break;

          //Host: Will Inject new rows for Online Bichitos. The first row is a tab for Offline Ones. 
          //Set "BID" on the injected row ID so it can be determined which bichito is each row later on 
          case "host":

            lower_ul.empty();
            lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"offline\" id=\""+htmlencode.htmlEncode(this.id)+"\">"+"Offlines"+"<span class=\"fa fa-chevron-down\"></span></a><ul class=\"nav child_menu\" style=\"display: block;\"></ul></li>");
            
            //Get Parent Implant RespTime, compare it with "lastchecked" to see if it is online
            for (i = 0; i < bichitos.length; i++){
              var row = bichitos[i];

              try{
              var infoJson = JSON.parse(row.info);
              } catch (e){
                console.log(e);
              }
              var bidentity = infoJson.mac.replace(/\n/g, '') + infoJson.hostname.replace(/\n/g, '');

              //Client side redirector organization
              if ((row.implantname == link.closest('.implantli').find('.implant').attr("id")) && (bidentity == link.attr("id")) && (row.status == "Online")){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"bichito\" id=\""+htmlencode.htmlEncode(row.bid)+"\">"+htmlencode.htmlEncode(row.bid)+"</a></li>");
              }
            }
            break;

          //Offline: Inject the rest of Offline, or previosly disconnected Bichtos  
          case "offline":

            lower_ul.empty();
            
            //Get Parent Implant RespTime, compare it with "lastchecked" to see if it is online
            for (i = 0; i < bichitos.length; i++){
              var row = bichitos[i];
 
              try{
              var infoJson = JSON.parse(row.info);
              } catch (e){
                console.log(e);
              }
              var bidentity = infoJson.mac.replace(/\n/g, '') + infoJson.hostname.replace(/\n/g, '');

              //Client side redirector organization
              if ((row.implantname == link.closest('.implantli').find('.implant').attr("id")) && (bidentity == link.attr("id")) && (row.status == "Offline")){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"bichito\" id=\""+htmlencode.htmlEncode(row.bid)+"\">"+htmlencode.htmlEncode(row.bid)+"</a></li>");
              }
            }
            break;

          //Bichito: Load main Bichito Component
          case "bichito":
            getBichitos();
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").empty();
            $(".STmain").load('./components/bichito/bichito.html')
            lower_ul.empty();
            break;

          // Vps,domains,Staging List/Menus, its functionality is similar to the one commented for Implants
          case "vpsList":
            getVps();
            lower_ul.empty();
            for (i = 0; i < vps.length; i++){
              var row = vps[i];
              lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"vps\" id=\""+htmlencode.htmlEncode(row.name)+"\">"+htmlencode.htmlEncode(row.name)+"</a></li>");
            }
            break;
          case "vps":
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/vps/vps.html')
            break;
          case "domainsList":
            getDomains();
            lower_ul.empty();
            for (i = 0; i < domains.length; i++){
              var row = domains[i];
              if (row.dtype != "gmail"){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"domain\" id=\""+htmlencode.htmlEncode(row.name)+"\">"+htmlencode.htmlEncode(row.name)+"</a></li>");
              }
            }
            break;
          case "saasList":
            getDomains();
            lower_ul.empty();
            for (i = 0; i < domains.length; i++){
              var row = domains[i];
              if (row.dtype == "gmail"){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"domain\" id=\""+htmlencode.htmlEncode(row.name)+"\">"+htmlencode.htmlEncode(row.name)+"</a></li>");
              }            }
            break;
          case "domain":
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/domain/domain.html')
            break;
          case "dropletList":
            getStagings();
            lower_ul.empty();
            for (i = 0; i < stagings.length; i++){
              var row = stagings[i];
              if (row.stype == "https_droplet_letsencrypt"){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"staging\" id=\""+htmlencode.htmlEncode(row.name)+"\">"+htmlencode.htmlEncode(row.name)+"</a></li>");
              }
            }
            break;
          case "hanlerList":
            getStagings();
            lower_ul.empty();
            for (i = 0; i < stagings.length; i++){
              var row = stagings[i];
              if (row.stype != "https_droplet_letsencrypt"){
                lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"staging\" id=\""+htmlencode.htmlEncode(row.name)+"\">"+htmlencode.htmlEncode(row.name)+"</a></li>");
              }
            }
            break;
          case "staging":
            getStagings();
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/staging/staging.html')
            break;
          case "reportsList":
            getReports();
            lower_ul.empty();
            for (i = 0; i < reports.length; i++){
              var row = reports[i];
              lower_ul.append("<li class=\"sub_menu\"><a href=\"#\" class=\"report\" id=\""+htmlencode.htmlEncode(row.name)+"\">"+htmlencode.htmlEncode(row.name)+"</a></li>");
            }
            break;
          case "report":
            $(".STmain").attr("id",link.attr("id"));
            $(".STmain").load('./components/report/report.html')
            break;


          //Creation Menus: Clicks on creation tabs, to load creation components
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
          case "addoperator":
            $(".STmain").load('./components/createforms/createOperator.html')
            break;
          default:
        }
    })
})









