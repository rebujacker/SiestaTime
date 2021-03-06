
// Loads Stagings for the Attacks
function loadFormData(button) {

  if (stagings != null){
    for (i = 0; i < stagings.length; i++){
      var row = stagings[i];
      $("#stagingOpt").append("<option>"+htmlencode.htmlEncode(row.name)+"</option>");
    }
  }

}

$(document).ready(function() {



//// Refresh on memory data and load it the sidetables for element creations
getStagings();
loadFormData();

//// Staging Creation Form: Parameters and function
  var name = $(".STmain").attr("id");
  $(".STmain").find(".element").text(htmlencode.htmlEncode(name));
  $("#delval").attr("value",htmlencode.htmlEncode(name));

  
  for (i = 0; i < implants.length; i++){
    if (implants[i].name == name){
      var implant = implants[i];
    }
  }

  var binumber = 0;
  for (i = 0; i < bichitos.length; i++){
    if (bichitos[i].implantname == name){
      binumber++;
    }
  }

  var idomains = [];
  var ivps = [];
  for (i = 0; i < redirectors.length; i++){
    if (redirectors[i].implantname == name){
      ivps.push(htmlencode.htmlEncode(redirectors[i].vpsname));
      idomains.push(htmlencode.htmlEncode(redirectors[i].domainame));
    }
  }

  var infoJson = JSON.parse(implant.modules)
  
  $("#network").text(htmlencode.htmlEncode(infoJson.coms.replace(/\n/g, '')));
  
  if (infoJson.persistence != undefined){
    $("#persistence").text(htmlencode.htmlEncode(infoJson.persistence.replace(/\n/g, '')));
  }else{
    $("#persistence").text("None");
  }

  $("#ivps").text(ivps);
  $("#idomains").text(idomains);
  $("#ibichitos").text(binumber);


//Used to change Forms for different staging types
$('#attacks').change(function(){

  switch($('#attacks').val()) {
    case 'dropImplant':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="attackparamsform">
    <div class="form-group">
        <label for="sel1">Target Operating System</label>
            <select class="form-control" id="os" name="stype">
                <option value="Linux">Linux</option>
                <option value="Windows">Windows</option>
                <option value="OSX">Mac OSX</option>
            </select>
    </div>

    <div class="form-group">
        <label for="sel1">Target Device Architecture</label>
            <select class="form-control" id="arch" name="stype">
                <option value="x32">x32</option>
                <option value="x64">x64</option>
            </select>
    </div>

    <div class="form-group">
      <label for="iname">Filename</label>
      <input type="text" class="form-control" id="filename" placeholder="goodboy.exe">
    </div>
    </form>   
    `);
      break;

    case 'hta':
      $("#params").empty();
      $("#params").append(`
    <form role="form" id="attackparamsform">
    <div class="form-group">
        <label for="sel1">Target Operating System</label>
            <select class="form-control" id="os" name="stype">
                <option value="Linux">Linux</option>
                <option value="Windows">Windows</option>
                <option value="OSX">Mac OSX</option>
            </select>
    </div>

    <div class="form-group">
        <label for="sel1">Target Device Architecture</label>
            <select class="form-control" id="arch" name="stype">
                <option value="x32">x32</option>
                <option value="x64">x64</option>
            </select>
    </div>

    <div class="form-group">
      <label for="iname">Filename</label>
      <input type="text" class="form-control" id="filename" placeholder="goodboy.exe">
    </div>
    </form>   
    `);
      break;
  }
});


  //Determine if elog is for hive,red or bichito, then load the logs for each case
  var name = $(".STmain").attr("id");
  $(".element").text(htmlencode.htmlEncode(name));
  

  
  $("#delval").attr("value",htmlencode.htmlEncode(name));

})

/* This Job will respect the following JSON Structure on "parameters":
type DeleteImplant struct{
    Name string `json:"name"`
}
*/

$("#submitdelimplant").on('click',function() {

  //Transform the array in one JSON STRING
  function objectifyForm(formArray) {
    var returnArray = {};
    for (var i = 0; i < formArray.length; i++){
        returnArray[formArray[i]['name']] = formArray[i]['value'];
      }
    return returnArray;
  }
  //Serialize form in the correct way

  var submitdelimplantJSON = objectifyForm($("#delimplantform").serializeArray());
  console.log(submitdelimplantJSON);

  //Create Job to send with two elements
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:"deleteImplant",time:"",status:"",result:"",parameters:"["+JSON.stringify(submitdelimplantJSON)+"]"};

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

/* This Job will respect the following JSON Structure on "parameters":
type DropImplant struct {
  Implant string   `json:"implant"`
    Staging string   `json:"staging"`
    Os string   `json:"os"`
    Arch string   `json:"arch"`
    Filename string   `json:"filename"`
}
*/

$("#submitattack").on('click',function(){

  var attack = $('#attacks').val();
  var implantName = $(".STmain").attr("id");
  //Create Job to send with two elements

  var createAttackJSON = {implant:implantName,staging:$("#stagingOpt").val(),os:$("#os").val(),arch:$("#arch").val(),filename:$("#filename").val()}
  var data = {cid:"",jid:"",pid:"Hive",chid:"None",job:attack,time:"",status:"",result:"",parameters:"["+JSON.stringify(createAttackJSON)+"]"};
 
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

/* Craft a Job with the following JSON Object towards client:

type ImplantObject struct {
    Name string   `json:"name"`
    OsName string   `json:"osname"`
    Arch string   `json:"arch"`
}
*/
$("#downloadImplant").on('click',function() {

  var implantName = $(".STmain").attr("id");

  //Create Job to send with two elements
  var data = {name:implantName,osname:$("#osD").val(),arch:$("#archD").val()};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/implant",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          console.log("Interact Sent!");

        }

    });

});

/* Craft a Job with the following JSON Object towards client:

type ImplantObject struct {
    Name string   `json:"name"`
    OsName string   `json:"osname"`
    Arch string   `json:"arch"`
}
*/
$("#downloadRedirector").on('click',function() {

  var implantName = $(".STmain").attr("id");

  //Create Job to send with two elements
  var data = {name:implantName,osname:"None",arch:"None"};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/redirector",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          console.log("Interact Sent!");

        }

    });

});


/* Craft a Job with the following JSON Object towards client:

type ShellcodeObject struct {
    Name string   `json:"name"`
    OsName string   `json:"osname"`
    Arch string   `json:"arch"`
    Format string   `json:"format"`
}
*/
$("#downloadShellcodeBin").on('click',function() {

  var implantName = $(".STmain").attr("id");

  //Create Job to send with two elements
  var data = {name:implantName,osname:$("#osD").val(),arch:$("#archD").val(),format:"binary"};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/shellcode",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          console.log("Interact Sent!");

        }

    });

});


/* Craft a Job with the following JSON Object towards client:

type ShellcodeObject struct {
    Name string   `json:"name"`
    OsName string   `json:"osname"`
    Arch string   `json:"arch"`
    Format string   `json:"format"`
}
*/
$("#downloadShellcodeHex").on('click',function() {

  var implantName = $(".STmain").attr("id");

  //Create Job to send with two elements
  var data = {name:implantName,osname:$("#osD").val(),arch:$("#archD").val(),format:"hex"};

  $.ajax({
        type: "POST",
        url: "http://127.0.0.1:8000/shellcode",
        data:  JSON.stringify(data),
        contentType: "application/json; charset=utf-8",
        dataType: "json",
        success: function (response){
          console.log("Interact Sent!");

        }

    });

});