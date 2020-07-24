
$(document).ready(function() {


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");
  
  $("#iname").text(htmlencode.htmlEncode(name));

  $("#delval").attr("value",htmlencode.htmlEncode(name));

  for (i = 0; i < domains.length; i++){
    if (domains[i].name == name){
      var domain = domains[i];
    }
  }

  $("#type").text(htmlencode.htmlEncode(domain.dtype));
  $("#domain").text(htmlencode.htmlEncode(domain.domain));
  $("#active").text(htmlencode.htmlEncode(domain.active));

})

/* This Job will respect the following JSON Structure on "parameters":
type DeleteDomain struct{
    Name string `json:"name"`
}
*/

$("#submitdeldomain").on('click',function() {

  //Transform the array in one JSON STRING
  function objectifyForm(formArray) {
    var returnArray = {};
    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    return returnArray;
  }

  //Serialize form in the correct way
  var submitdeldomainJSON = objectifyForm($("#deldomainform").serializeArray());

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"deleteDomain",time:"",status:"",result:"",parameters:"["+JSON.stringify(submitdeldomainJSON)+"]"};

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
