
$(document).ready(function() {


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");

  $("#iname").text(htmlencode.htmlEncode(name));
  $("#element").text(htmlencode.htmlEncode(name));
  $("#delval").attr("value",htmlencode.htmlEncode(name));

  for (i = 0; i < vps.length; i++){
    if (vps[i].name == name){
      var vpsi = vps[i];
    }
  }

  $("#type").text(htmlencode.htmlEncode(vpsi.vtype));


})

/*
type DeleteVps struct{
    Name string `json:"name"`
}
*/
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