$(document).ready(function () {
  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  let tableParam = urlParams.get("table");
  const table = tableParam ? tableParam : "Charts";

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
    if(type == "Post") {
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
      titleText = "PUBLISHED POSTS PER MONTH";
    } else if (type = "Category") {
      data.DataChart.forEach((element) => (dataC.push(element.Critere)));
      var lab = dataC;
      titleText = "PUBLISHED POSTS PER CATEGORIES";
    }
    data.DataChart.forEach((element) => (dataC[element.Critere] = element.Count));
    $("#chartCont").html(
      '<canvas id="myChart" width="825" height="400"></canvas>'
    );
    var chart = document.getElementById("myChart");
    var myMonthlyChart = new Chart(chart, {
      type: "bar",
      data: {
        labels: lab,
        datasets: [
          {
            stack: 0,
            data: dataC,
            backgroundColor: ["rgba(0,243,255,1)"],
            borderColor: ["rgba(2,216,227, 1)"],
            borderWidth: 1,
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
        
        tooltips: {},
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
          if (tableAdmin != null) {
            tableAdmin.clear().rows.add(x.Post).draw(false);
          } else {
            tableAdmin = $("#tableAdmin").DataTable({
              columns: [
                {
                  title: "ID",
                  data: "ID",
                },
                {
                  title: "TITLE",
                  data: "Title",
                },
                {
                  title: "AUTHOR",
                  data: "AuthorID",
                },
                {
                  title: "CONTENT",
                  data: "Content",
                  render: function (data) {
                    return "";
                  },
                },
                {
                  title: "CATEGORY",
                  data: "CategoryID",
                },
                {
                  title: "DATE",
                  data: "Date",
                  render: function (data) {
                    return data + "fuck";
                  },
                },
                {
                  title: "STATE",
                  data: "State",
                },
                {
                  title: "",
                  data: "ID",
                  render: function (data) {
                    return data + "IconEdit";
                  },
                },
              ],
              data: x.Post,
              order: [[0, "asc"]],
            });
          }
        } else {
          // Si onglet actif = Charts
          countElements(x);
          console.log(x);
        }
      });
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
