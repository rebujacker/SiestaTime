
$(document).ready(function() {


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");
  //console.log(name)
  $(".element").text(htmlencode.htmlEncode(name))
  //$("#delval").value = name
  $(".element").attr("value",htmlencode.htmlEncode(name));

})

$("#submitdownloadreport").on('click',function() {

  var reportName = $(".element").attr("value");

  //Create Job to send with two elements
  var data = {name:reportName};
  //data.push();
  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/report",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          //console.log("Interact Sent!");

        }

    });

});

/*

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
  //console.log(submitdeldomainJSON);



  function delDomain() {

    return $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/deldomain",
        data:  JSON.stringify(submitdeldomainJSON),
        contentType: "application/json; charset=utf-8",
        dataType: "json"         
    });
  }

  $.when(delDomain()).done(function(){

    //Toast Domain Created/Error
  });

});

*/