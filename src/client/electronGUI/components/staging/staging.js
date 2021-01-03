$(document).ready(function() {


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");

  $("#iname").text(htmlencode.htmlEncode(name));
  $("#delval").attr("value",htmlencode.htmlEncode(name));

  for (i = 0; i < stagings.length; i++){
    if (stagings[i].name == name){
      var staging = stagings[i];
    }
  }

  $("#type").text(htmlencode.htmlEncode(staging.stype));
  $("#vps").text(htmlencode.htmlEncode(staging.vpsname));
  $("#domain").text(htmlencode.htmlEncode(staging.domainame));

});


$('#interactType').change(function(){

  if ($('#interactType').val() == 'socks5'){
    $("#interactParams").empty();
    $("#interactParams").append(`
    
      <div class="form-group">
        <label for="timeouts">Operator Port (127.0.0.1)</label>
        <input type="text" class="form-control" id="socks5portOpt" name="socks5port" placeholder="socks5 port...">
      </div>

      <button type="button" class="btn btn-primary" id="interact">Open Socks5 in my device</button>
    `);

  }else if ($('#interactType').val() == 'kilssh'){
    $("#interactParams").empty();
    $("#interactParams").append(`
    
      <button type="button" class="btn btn-primary" id="interact">Kill Active SSH Sessions</button>
    `);  
  
}else if ($('#interactType').val() == 'fullinteractive'){

  }else{
    $("#interactParams").empty();
    $("#interactParams").append(`
    
      <button type="button" class="btn btn-primary" id="interact">Interact</button>
    `); 
  }

});


/* Craft a Job with the following JSON Object towards client:

type InteractObject struct {
    StagingName string   `json:"staging"`
    Handler string   `json:"handler"`
    VpsName string   `json:"vpsname"`
    TunnelPort string   `json:"tunnelport"`
    Socks5Port string   `json:"socks5port"`
}
*/
$("#interactParams").on('click','#interact',function(){

  var stagingName = $(".STmain").attr("id");
  var vpsName = "";
  var tunnelPort = "";
  var handlerN = "";
  var socks5Port = $('#socks5portOpt').val();;

  for (i = 0; i < stagings.length; i++){
    if(stagings[i].name == stagingName) {
      vpsName = stagings[i].vpsname
      tunnelPort = stagings[i].tunnelport
      switch(stagings[i].stype) {
        case "https_droplet_letsencrypt":
          handlerN = "droplet";
          break;
        case "https_msft_letsencrypt":
          handlerN = "msfconsole";
          break;
        case "https_empire_letsencrypt":
          handlerN = "empire";
          break;
        case "ssh_rev_shell":
          if ($('#interactType').val() == 'fullinteractive'){
            handlerN = "ssh";
            break;
          }else if ($('#interactType').val() == 'socks5'){

          handlerN = "socks5";
          break;          
        }else{
          handlerN = "killssh";
          break;            
        }
      
      }

    }
  }

  //Create Job to send with two elements
  var data = {staging:stagingName,handler:handlerN,vpsname:vpsName,tunnelport:tunnelPort,socks5port:socks5Port};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/interact",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){

        }

    });
});


/*
type DeleteStaging struct{
    Name string `json:"name"`
}
*/

$("#killsshSession").on('click',function() {

  var stagingName = $(".STmain").attr("id");
  var vpsName = "";
  var tunnelPort = "";
  var handlerN = "killSSH";
  var socks5Port = $('#socks5portOpt').val();;


  for (i = 0; i < stagings.length; i++){
    if(stagings[i].name == stagingName) {
      vpsName = stagings[i].vpsname
      tunnelPort = stagings[i].tunnelport
      switch(stagings[i].stype) {
        case "https_droplet_letsencrypt":
          handlerN = "droplet";
          break;
        case "https_msft_letsencrypt":
          handlerN = "msfconsole";
          break;
        case "https_empire_letsencrypt":
          handlerN = "empire";
          break;
        case "ssh_rev_shell":
          if ($('#interactType').val() == 'fullinteractive'){
            handlerN = "ssh";
            break;
          }

          handlerN = "socks5";
          break;          

      }

    }
  }


  //Create Job to send with two elements
  var data = {staging:stagingName,handler:handlerN,vpsname:vpsName,tunnelport:tunnelPort,socks5port:socks5Port};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/interact",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){

        }

    });
});


/*
type DeleteStaging struct{
    Name string `json:"name"`
}
*/

$("#submitdelstaging").on('click',function() {

  //Transform the array in one JSON STRING
  function objectifyForm(formArray) {
    var returnArray = {};
    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    return returnArray;
  }
  //Serialize form in the correct way

  var submitdelstagingJSON = objectifyForm($("#delstagingform").serializeArray());
    
  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"deleteStaging",time:"",status:"",result:"",parameters:"["+JSON.stringify(submitdelstagingJSON)+"]"};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/job",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){

          if (response != null){

            return
          }
        }

    });
});