$(document).ready(function() {


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");
  //console.log(name)
  $("#iname").text(htmlencode.htmlEncode(name));
  //$("#delval").value = name
  $("#delval").attr("value",htmlencode.htmlEncode(name));

  for (i = 0; i < stagings.length; i++){
    if (stagings[i].name == name){
      var staging = stagings[i];
    }
  }

  $("#type").text(htmlencode.htmlEncode(staging.stype));
  $("#vps").text(htmlencode.htmlEncode(staging.vpsname));
  $("#domain").text(htmlencode.htmlEncode(staging.domainame));

})


$("#interact").on('click',function() {

  var stagingName = $(".STmain").attr("id");
  var vpsName = "";
  var tunnelPort = "";
  var handlerN = "";

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
          handlerN = "ssh";
          break;

      }

    }
  }

  //Create Job to send with two elements
  var data = {staging:stagingName,handler:handlerN,vpsname:vpsName,tunnelport:tunnelPort};
  //data.push();
  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/interact",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          //console.log("Interact Sent!");

        }

    });
});


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
  //console.log(submitdelstagingJSON);
    
  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"deleteStaging",time:"",status:"",result:"",parameters:"["+JSON.stringify(submitdelstagingJSON)+"]"};
  //data.push();
  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/job",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          //console.log("Response Job:"+response[0].jid);
          if (response != null){
            //console.log("Response Job:"+response[0].jid);
            return
          }
        }

    });
});