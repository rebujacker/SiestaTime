
$(document).ready(function() {


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");
  //console.log(name)
  $("#iname").text(htmlencode.htmlEncode(name));
  $("#element").text(htmlencode.htmlEncode(name));
  //$("#delval").value = name
  $("#delval").attr("value",htmlencode.htmlEncode(name));

  for (i = 0; i < vps.length; i++){
    if (vps[i].name == name){
      var vpsi = vps[i];
    }
  }

  $("#type").text(htmlencode.htmlEncode(vpsi.vtype));


})

$("#submitdelvps").on('click',function() {

  //Transform the array in one JSON STRING
  function objectifyForm(formArray) {
    var returnArray = {};
    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    return returnArray;
  }
  //Serialize form in the correct way

  var submitdelvpsJSON = objectifyForm($("#delvpsform").serializeArray());
  //console.log(submitdelvpsJSON);
  
  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"deleteVPS",time:"",status:"",result:"",parameters:"["+JSON.stringify(submitdelvpsJSON)+"]"};
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

/*

$("#submitdelvps").on('click',function() {

  //Transform the array in one JSON STRING
  function objectifyForm(formArray) {
    var returnArray = {};
    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    return returnArray;
  }
  //Serialize form in the correct way

  var submitdelvpsJSON = objectifyForm($("#delvpsform").serializeArray());
  //console.log(submitdelvpsJSON);



  function delImplant() {

    return $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/delvps",
        data:  JSON.stringify(submitdelvpsJSON),
        contentType: "application/json; charset=utf-8",
        dataType: "json"         
    });
  }

  $.when(delImplant()).done(function(){

    //Toast Implant Created/Error
  });

});
*/