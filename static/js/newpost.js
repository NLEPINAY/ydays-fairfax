// Lorsque la page est totalement chargée...
$(document).ready(function () {
  // Centre le livre au milieu de la page
  $(document).on("click", "#page-1", function () {
    $("#main").css("left", "calc(50vw - (var(--baseline) * 60)/1.35)");
  });

  $(document).on("click", "#page-2", function () {
    $("#main").css("left", "calc(50vw - (var(--baseline) * 60)/2)");
  });

  // Initialisation tinymce (éditeur de texte WYSIWYG)
  tinymce.init({
    menubar: false,
    language: "fr_FR",
    selector: "#textarea-content",
    plugins:
      "autoresize image template textcolor hr importcss link noneditable",
    //"lists advlist autoresize charmap emoticons media image template preview textcolor hr importcss link noneditable table help",
    toolbar:
      "fontselect fontsizeselect | bold italic forecolor" +
      "alignleft aligncenter alignright alignjustify | " +
      "imageupload image",
    target_list: [
      { title: "Téléchargement", value: "_self" },
      { title: "Nouvel onglet", value: "_blank" },
    ],
    media_live_embeds: true,
    browser_spellcheck: true,
    spellchecker_language: "fr-FR",
    noneditable_noneditable_class: "mceNonEditable",
    branding: false,
    mobile: {
      menubar: true,
    },
    file_picker_types: "file image media",
    images_upload_handler: function (blobInfo, success, failure) {
      var xhr, formData;
      xhr = new XMLHttpRequest();
      xhr.withCredentials = false;
      xhr.open("POST", "./private/controller/traitement-ajax.php");
      xhr.onload = function () {
        var json;
        if (xhr.status != 200) {
          failure("HTTP Error: " + xhr.status);
          return;
        }
        json = JSON.parse(xhr.responseText);

        if (!json || typeof json.location != "string") {
          failure("Invalid JSON: " + xhr.responseText);
          return;
        }
        success(json.location);
      };
      formData = new FormData();
      formData.append("file", blobInfo.blob(), blobInfo.filename());
      xhr.send(formData);
    },

    file_picker_callback: function (cb, value, meta) {
      var input = document.createElement("input");
      input.setAttribute("type", "file");
      input.setAttribute("accept", "image/* audio/* video/*");
      input.onchange = function () {
        var file = this.files[0];

        var reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = function () {
          var id = "blobid" + new Date().getTime();
          var blobCache = tinymce.activeEditor.editorUpload.blobCache;
          var base64 = reader.result.split(",")[1];
          var blobInfo = blobCache.create(id, file, base64);
          blobCache.add(blobInfo);

          // call the callback and populate the Title field with the file name
          cb(blobInfo.blobUri(), {
            title: file.name,
          });
        };
      };

      input.click();
    },
  });

  $(document).on("focusin", function (e) {
    if ($(e.target).closest(".mce-window").length) {
      e.stopImmediatePropagation();
    }
  });
});
