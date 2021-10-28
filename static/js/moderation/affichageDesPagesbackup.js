const formInputs = {
  Post: {
    Title: {
      type: "text",
    },    
    Category: {
      type: "select",
      table: "Category",
    },
    Content: {
      type: "textarea",
    },
  },
  Comment: {
    Content: {
      type: "textarea",
    },
  },
  User: {},
  Category: {},
};

$(document).ready(function () {
  
  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  let tableParam = urlParams.get("table");
  var table = tableParam ? tableParam : "Charts";

  $(document).on("click", ".catBtn", function () {
    $(".catBtn").removeClass("activeCat");
    $(this).addClass("activeCat");
    if ($(this).attr("data-table") == "Charts") {
      $("#table").addClass("hidden");
      $("#cardsContainer").removeClass("hidden");
    } else {
      $("#table").removeClass("hidden");
      $("#cardsContainer").addClass("hidden");
    }
    initDatatable($(this).attr("data-table"));
  });

  // Affichage des différentes stats
  $(document).on("click", "#topCont .card", function () {
    $("#topCont .card").removeClass("activeChart");
    $(this).addClass("activeChart");
    const type = $(this).attr("data-chart");

    var params = new Object();
    params.action = "getStats";
    params.what = type;
    fetch("/fetching", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(params),
    })
      .then((x) => x.json())
      .then((x) => {
        generateChart(x, type);
      });
  });

  //Initialisation menu actif
  $("li[data-table='" + table + "']").click();

  function countElements(response) {
    const types = [
      ["postsCount", response.Post[0].Count],
      ["categoriesCount", response.Category[0].Count],
      ["usersCount", response.User[0].Count],
      ["commentsCount", response.Comment[0].Count],
    ];
    types.forEach((element) =>
      new CountUp(element[0], 0, element[1], 0, 2.5).start()
    );
  }

  /**
   * Génération des graphiques avec Chart.js
   * Documentation : https://www.chartjs.org/docs/latest/getting-started/
   */
  async function generateChart(data, type) {
    console.log("data:", data);
    var dataC = [];
    var titleText = "";
    var typeChart = "";
    if (type == "Post" || type == "User") {
      dataC = ["0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0", "0"];
      var lab = [
        "January",
        "February",
        "March",
        "April",
        "May",
        "June",
        "July",
        "August",
        "September",
        "October",
        "November",
        "December",
      ];
      typeChart = "line";
      titleText =
        type == "Post"
          ? "PUBLISHED POSTS PER MONTH"
          : "USERS REGISTRATION PER MONTH";
    } else if (type == "Category") {
      data.DataChart.forEach((element) => dataC.push(element.Critere));
      var lab = dataC;
      titleText = "PUBLISHED POSTS PER CATEGORIES";
      typeChart = "bar";
    }
    data.DataChart.forEach(
      (element) => (dataC[element.Critere - 1] = element.Count)
    );
    $("#chartCont").html(
      '<canvas id="myChart" width="825" height="400"></canvas>'
    );
    var chart = document.getElementById("myChart");
    var myMonthlyChart = new Chart(chart, {
      type: typeChart,
      data: {
        labels: lab,
        datasets: [
          {
            //stack: 0,
            data: dataC,
            backgroundColor: ["rgba(0,243,255,1)"],
            borderColor: ["rgba(5,170,223,1)"],
            borderWidth: 1,
            pointBackgroundColor: "rgba(92, 145, 249, 1)",
            lineTension: 0.4,
            fill: true,
          },
        ],
      },
      options: {
        plugins: {
          legend: {
            display: false,
          },
          title: {
            display: true,
            text: titleText,
            fontColor: "rgba(0,243,255,1)",
            font: {
              size: 20,
            },
          },
        },

        tooltips: { mode: "point" },
        scales: {
          xAxes: [
            {
              stacked: true,
              display: true,
              ticks: {
                beginAtZero: true,
                min: 0,
                suggestedMin: 0,
              },
            },
          ],
          yAxes: [
            {
              stacked: true,
              display: true,
              ticks: {
                beginAtZero: true,
                min: 0,
                suggestedMin: 0,
              },
            },
          ],
        },
      },
    });
  }

  var tableAdmin = null;

  function initDatatable(table) {
    var params = new Object();
    params.action = "get";
    params.table = table;
    fetch("/fetching", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(params),
    })
      .then((x) => x.json())
      .then((x) => {
        if (table != "Charts") {
          const dataArray = [];
          var dataAttributes = "";
          const columnsToExclude = [
            "Liked",
            "Disliked",
            "Image",
            "Gif",
            "Password",
            "SecretQuestion",
            "SecretAnswer",
            "Count",
            "AuthorID",
            "Reason",
            "Content",
            "CategoryID",
          ];
          if(table == "Post") {
            // Ajout du nombre de likes / dislikes pour chaque post
            x[table].forEach((ind) => { // Boucle sur x.Post
              ind['Likes'] = 0;    // Valeurs par défaut
              ind['Dislikes'] = 0;
              x['CountLike'].forEach((index) => { // Boucle sur x.CountLike
                if (ind['ID'] == index['PostId']) {
                  ind['Likes'] = index['CountLikes'] + '<i class="fas fa-heart ml-2 text-danger"></i>';
                  ind['Dislikes'] = index['CountDislikes'] + '<i class="fas fa-heart-broken ml-2"></i>';
                }
              })
            })           
          }
          Object.entries(x[table][0]).forEach(([key, value]) => {
            dataAttributes += "data-" + key.toLowerCase() + '="' + value + '" ';
            if (!columnsToExclude.includes(key)) {
              var column = {};
              column.title = key
                .replace("ID", "Id")
                .replace(/([A-Z])/g, " $1")
                .toUpperCase(); //.replace(/_/g, " ") => pour transformer du snake case;
              column.data = key;
              // Paramètres spécifiques
              if (key == "Date") {
                column.render = function (data) {
                  return convertirDate(data);
                };
              } else if (key == "Avatar") {
                column.render = function (data) {
                  return (
                    `<div class="avatar" style="background-image:URL('.` +
                    data +
                    `');"></div>`
                  );
                };
              } else if (key == "House") {
                column.render = function (data) {
                  return (
                    `<div class="house infoLien" style="background-image:URL('..` +
                    data.Image +
                    `');"><span>` +
                    data.Name +
                    `</span></div>`
                  );
                };
              } else if (key == "Author") {
                column.render = function (data) {
                  return data.Username;
                };
              } else if (key == "Category") {
                column.render = function (data) {
                  return data.Name;
                };
              } else if (key.includes("State")) {
                column.render = function (data) {
                  var result =
                    data == 1
                      ? '<i class="fas fa-eye"></i>'
                      : '<i class="fas fa-eye-slash"></i>';
                  return result;
                };
              }
              dataArray.push(column);
            }
          });
          // Paramètres globaux
          dataArray.push({
            class: "action",
            title: "ACTIONS",
            orderable: false,
            render: function (data, type, row, meta) {
              return (
                '<div class="btn btn-outline-blue editLink mr-3" data-id="' +
                row.ID +
                '" data-table-type="' +
                table +
                '" ' +
                dataAttributes +
                ' data-toggle="modal" data-target="#updateModal"><i class="far fa-edit"></i></div><div class="btn btn-outline-danger deleteLink" data-id="' +
                row.ID +
                '"><i class="far fa-trash-alt"></i></div>'
              );
            },
          });
          if (tableAdmin != null) {
            tableAdmin.destroy();
            $("#tableAdmin").empty();
          }
          tableAdmin = $("#tableAdmin").DataTable({
            columns: dataArray,
            data: x[table],
            scrollX: true,
            autoWidth: false,
            order: [[0, "desc"]],
          });
        } else {
          // Si onglet actif = Charts
          countElements(x);
          $("div[data-chart='Post']").click();
        }
      });
  }

  function convertirDate(date) {
    var monthName = [
      "Janvier",
      "Février",
      "Mars",
      "Avril",
      "Mai",
      "Juin",
      "Juillet",
      "Août",
      "Septembre",
      "Octobre",
      "Novembre",
      "Décembre",
    ];
    var dayName = [
      "Dimanche",
      "Lundi",
      "Mardi",
      "Mercredi",
      "Jeudi",
      "Vendredi",
      "Samedi",
    ];

    var maDate = new Date(date);
    var jour = maDate.getDay(); //Jour
    var njour = maDate.getDate(); //Numéro du jour
    var mois = maDate.getMonth(); //Mois (commence à 0, donc +1)
    var annee = maDate.getFullYear(); //Année sur 2 chiffres ou getFullYear sur 4

    var resultDate = njour + " " + monthName[mois] + " " + annee;
    return resultDate;
  }

  initDatatable(table);

  // Modal Update
  $(document).on("click", ".editLink", function () {
    var params = new Object();
    const tableToUpdate = $(this).attr("data-table-type");
    params.action = "getForUpdate";
    params.table = tableToUpdate;
    params.id = $(this).attr("data-id");
    fetch("/fetching", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(params),
    })
      .then((x) => x.json())
      .then((x) => {
        
        var formContent = "";
        var textarea = "";
        var textareaContent = "";
        Object.entries(formInputs[tableToUpdate]).forEach(([input, data]) => {
            if (formInputs[`${tableToUpdate}`][`${input}`].type == "textarea") {
              textareaContent = x[`${tableToUpdate}`][0][`${input}`];
              textarea +=
              `<div class="form-group w-100">
              <label for="` +
              input +
              `">` +
              input +
              `</label>
              <textarea class="form-control" name="` +
                input +
                `" id="textareaContent">` +
                /*x[`${tableToUpdate}`][0][`${input}`] +
                */`</textarea></div>`;
            } else if (formInputs[`${tableToUpdate}`][`${input}`].type == "select") {
              var options = "";
              var tableSelect = formInputs[`${tableToUpdate}`][`${input}`].table;
              x.Category.forEach((element) => {
                options += `<option value="`+element.Name+`">`+element.Name+`</option>`;
              })
              formContent +=
              `<div class="form-group">
              <label for="` +
              input +
              `">` +
              input +
              `</label>
              <select class="form-control" name="`+input+`">
              `+options+`
              </select>
              </div>`
            } else {
              formContent +=
                `<div class="form-group">
                <label for="` +
                input +
                `">` +
                input +
                `</label>
                <input class="form-control" type="` +
                formInputs[`${tableToUpdate}`][`${input}`].type +
                `" name="` +
                input +
                `" value="` +
                x[`${tableToUpdate}`][0][`${input}`] +
                `"></div>`;
            }
        })

        formContent += textarea;
        $("#updateForm").html(formContent);
      
          // S'il y a un textarea, on initialise l'éditeur
          if(textarea) {
            // Initialisation tinymce (éditeur de texte WYSIWYG)
            tinymce.init({
              menubar: false,
              language: "fr_FR",
              selector: "#textareaContent",
              plugins:
                "autoresize image template textcolor hr importcss link noneditable spellchecker",
              //"lists advlist autoresize charmap emoticons media image template preview textcolor hr importcss link noneditable table help",
              toolbar:
                "fontselect fontsizeselect | bold italic forecolor" +
                "alignleft aligncenter alignright alignjustify | " +
                "imageupload image" + "| spellchecker",
              target_list: [
                { title: "Téléchargement", value: "_self" },
                { title: "Nouvel onglet", value: "_blank" },
              ],
              media_live_embeds: true,
              browser_spellcheck: true,
              spellchecker_languages: 'English=en,Danish=da,Dutch=nl,Finnish=fi,French=fr_FR,' +
          'German=de,Italian=it,Polish=pl,Portuguese=pt_BR,Spanish=es,Swedish=sv',
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
            }).then(tinymce.get("textareaContent").setContent(textareaContent));
            console.log("textarea-content: ",textareaContent);
            console.log("textarea: ",$("#textareaContent"));
            //tinymce.get("textareaContent").setContent(textareaContent);

            $(document).on("focusin", function (e) {
              if ($(e.target).closest(".mce-window").length) {
                e.stopImmediatePropagation();
              }
            });
          }
      });

    $(".modal-title").html(
      "UPDATE " + tableToUpdate.toUpperCase() + " N°" + $(this).attr("data-id")
    );

  });

  //Initialisation de dataTable avec des paramètres personnalisés
  if ($.fn.dataTable) {
    $.extend($.fn.dataTable.defaults, {
      language: {
        sEmptyTable: "Aucune donnée disponible dans le tableau",
        sInfo: "Affichage de l'élément _START_ à _END_ sur _TOTAL_ éléments",
        sInfoEmpty: "Affichage de l'élément 0 à 0 sur 0 élément",
        sInfoFiltered: "(filtré à partir de _MAX_ éléments au total)",
        sInfoPostFix: "",
        sInfoThousands: ",",
        sLengthMenu: "Afficher _MENU_ éléments",
        sLoadingRecords: "Chargement...",
        sProcessing: "Traitement...",
        sSearch: "Rechercher :",
        sZeroRecords: "Aucun élément correspondant trouvé",
        oPaginate: {
          sFirst: "Premier",
          sLast: "Dernier",
          sNext: "Suivant",
          sPrevious: "Précédent",
        },
        oAria: {
          sSortAscending: ": activer pour trier la colonne par ordre croissant",
          sSortDescending:
            ": activer pour trier la colonne par ordre décroissant",
        },
        select: {
          rows: {
            _: "%d lignes sélectionnées",
            0: "Aucune ligne sélectionnée",
            1: "1 ligne sélectionnée",
          },
        },
      },
    });
  }
});
