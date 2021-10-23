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
    const type = $(this).attr('data-chart')

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
          generateChart(x,type);
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

    console.log("data:",data);
    var dataC = [];
    var titleText = "";
    var typeChart = "";
    if((type == "Post")||(type == "User")) {
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
      titleText = type == "Post" ? "PUBLISHED POSTS PER MONTH" : "USERS REGISTRATION PER MONTH";
    } else if (type == "Category") {
      data.DataChart.forEach((element) => (dataC.push(element.Critere)));
      var lab = dataC;
      titleText = "PUBLISHED POSTS PER CATEGORIES";
      typeChart = "bar";
    }
    data.DataChart.forEach((element) => (dataC[element.Critere-1] = element.Count));
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
            pointBackgroundColor: 'rgba(92, 145, 249, 1)',
            lineTension: 0.4,  
            fill:true,
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
            size: 20
        }
        },
      },
        
        tooltips: {mode: 'point'},
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
          const columnsToExclude = ['Liked', 'Disliked','Image','Gif','Password','SecretQuestion','SecretAnswer','Count','Author','Reason','Content'];
          Object.entries(x[table][0]).forEach(([key, value]) => {
            if(!columnsToExclude.includes(key)) {
              var column = {};
                column.title = key.toUpperCase().replace(/_/g, " ");
                column.data = key;  
                // Paramètres spécifiques   
                if(key == "Date") {
                  column.render = function (data) {
                    return convertirDate(data);
                  }
                } else if (key == "Avatar") {
                  column.render = function (data) {
                    return `<div class="avatar" style="background-image:URL('.`+data+`');"></div>`;
                  }
                } else if (key == "House") {
                  column.render = function (data) {
                    console.log(data);
                    return `<div class="house infoLien" style="background-image:URL('..`+data.Image+`');"><span>`+data.Name+`</span></div>`;
                  }
                }                
              dataArray.push(column);
            }
          });
          // Paramètres globaux
          dataArray.push(
            {
              class: "action",
              title: "ACTIONS",
              orderable: false,
              render: function (data, type, row, meta) {
                return (
                  '<a href="#" class="infoLien editLink" data-id="' +
                  row.ID +
                  '" data-toggle="modal" data-target="#updateModal"><i class="far fa-edit"></i><span>Update</span></a><a href="#" class="infoLien deleteLink" data-id="' +
                  row.ID +
                  '"><i class="far fa-trash-alt"></i><span>Delete</span></a>'
                );
              },
            }
          );
            if (tableAdmin != null) {
              tableAdmin.destroy();
              $("#tableAdmin").empty();
            }
              tableAdmin = $("#tableAdmin").DataTable({
                columns: dataArray,
                data: x[table],
                "scrollX":true,
                "autoWidth": false,
                order: [[0, "desc"]],
              });
            
        } else {
          // Si onglet actif = Charts
          countElements(x);
          $("div[data-chart='Post']").click();
          console.log(x);
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
