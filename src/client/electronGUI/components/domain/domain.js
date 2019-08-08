
$(document).ready(function() {


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");
  //console.log(name)
  $("#iname").text(name);
  //$("#delval").value = name
  $("#delval").attr("value",name);

  for (i = 0; i < domains.length; i++){
    if (domains[i].name == name){
      var domain = domains[i];
    }
  }

  $("#type").text(domain.dtype);
  $("#domain").text(domain.domain);
  $("#active").text(domain.active);

})

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
  console.log(submitdeldomainJSON);

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"deleteDomain",time:"",status:"",result:"",parameters:"["+JSON.stringify(submitdeldomainJSON)+"]"};
  //data.push();
  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/job",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          console.log("Response Job:"+response[0].jid);
          if (response != null){
            console.log("Response Job:"+response[0].jid);
            return
          }
        }

    });

});
