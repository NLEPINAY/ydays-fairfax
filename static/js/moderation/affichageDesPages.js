$(document).ready(function () {
  const queryString = window.location.search;
  const urlParams = new URLSearchParams(queryString);
  let tableParam = urlParams.get("table");
  const table = tableParam ? tableParam : "posts";
  //history.pushState({}, null, "/moderation&table=" + table);

  $(document).on("click", ".catBtn", function () {
    initDatatable($(this).val());
  });

  var tableAdmin = null;

  function initDatatable(table) {
    // $.ajax({
    //   url: "/fetching", // on appelle le fichier de traitement
    //   type: "POST",
    //   //dataType: "json",
    //   data: {
    //     //param: table,
    //     action:"get",
    //   },
    //   contentType: false,
    //   processData: false,
    //   success: function (data) {
    //     if (tableAdmin != null) {
    //       tableAdmin.clear().rows.add(data).draw(false);
    //     } else {
    //       tableAdmin = $("#tableAdmin").DataTable({
    //         columns: [
    //           {
    //             title: "ID",
    //             data: "id",
    //           },
    //           {
    //             title: "TITLE",
    //             data: "titre",
    //           },
    //           {
    //             title: "AUTHOR",
    //             data: "titre",
    //           },
    //           {
    //             title: "CONTENT",
    //             data: "titre",
    //           },
    //           {
    //             title: "CATEGORY",
    //             data: "titre",
    //           },
    //           {
    //             title: "DATE",
    //             data: "titre",
    //           },
    //           {
    //             title: "STATE",
    //             data: "titre",
    //           },
    //         ],
    //         data: data,
    //         order: [[0, "desc"]],
    //       });
    //     }
    //   },
    //   error: function (data) {
    //     console.log("error");
    //     console.log(data);
    //   },
    // });
    var params = new Object()
    params.action = "get"
    fetch("/fetching", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                Accept: "application/json"
            },
            body: JSON.stringify(params)
        }).then(x => x.json()).then(x=>{
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
              },
              {
                title: "CATEGORY",
                data: "CategoryID",
              },
              {
                title: "DATE",
                data: "Date",
                render: function (data) {return data+"fuck"}
              },
              {
                title: "STATE",
                data: "State",
              },
              {
                title : "",
                data: "ID",
                render: function (data) {return data+"IconEdit"}
              },
            ],
            data: x.Post,
            order: [[0, "asc"]],
          });
        }
    
  })
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

/*
//Affiche seulement le tableau post et cache les autres. Chaque fonction est casi pareil, je cache les autres table SI elle sont visible et j'affiche celle si si elle ne l'est pas.
function showPosts() {
  if (!commentTable.className.includes("Invisible")) {
    commentTable.classList.add("Invisible");
    commentTable.classList.remove("Visible");
  }
  if (postTable.className.includes("Invisible")) {
    removeSelected();
    removeVisibility("postOption");
    postTable.classList.remove("Invisible");
    postTable.classList.add("Visible");
    if (actualTable != "posts") {
      actualTable = "posts";
      if (urlParams.get("page") == null) {
        history.pushState({}, null, "/moderation&table=posts");
      } else {
        history.pushState(
          {},
          null,
          "/moderation?page=" + urlParams.get("page") + "&table=posts"
        );
      }
    }
  }
  if (!userTable.className.includes("Invisible")) {
    userTable.classList.add("Invisible");
    userTable.classList.remove("Visible");
  }
  if (!categoryTable.className.includes("Invisible")) {
    categoryTable.classList.add("Invisible");
    categoryTable.classList.remove("Visible");
  }
}
//Affiche seulement le tableau comments et cache les autres. Chaque fonction est casi pareil, je cache les autres table SI elle sont visible et j'affiche celle si si elle ne l'est pas.
function showComments() {
  if (commentTable.className.includes("Invisible")) {
    removeSelected();
    removeVisibility("commentOption");
    commentTable.classList.remove("Invisible");
    commentTable.classList.add("Visible");
    if (actualTable != "comments") {
      actualTable = "comments";
      if (urlParams.get("page") == null) {
        history.pushState({}, null, "/moderation&table=comments");
      } else {
        history.pushState(
          {},
          null,
          "/moderation?page=" + urlParams.get("page") + "&table=comments"
        );
      }
    }
  }
  if (!postTable.className.includes("Invisible")) {
    postTable.classList.add("Invisible");
    postTable.classList.remove("Visible");
  }
  if (!userTable.className.includes("Invisible")) {
    userTable.classList.add("Invisible");
    userTable.classList.remove("Visible");
  }
  if (!categoryTable.className.includes("Invisible")) {
    categoryTable.classList.add("Invisible");
    categoryTable.classList.remove("Visible");
  }
}
//Affiche seulement le tableau users et cache les autres. Chaque fonction est casi pareil, je cache les autres table SI elle sont visible et j'affiche celle si si elle ne l'est pas.
function showUsers() {
  if (!commentTable.className.includes("Invisible")) {
    commentTable.classList.add("Invisible");
    commentTable.classList.remove("Visible");
  }
  if (!postTable.className.includes("Invisible")) {
    postTable.classList.add("Invisible");
    postTable.classList.remove("Visible");
  }
  if (userTable.className.includes("Invisible")) {
    removeSelected();
    removeVisibility("userOption");
    userTable.classList.remove("Invisible");
    userTable.classList.add("Visible");
    if (actualTable != "users") {
      actualTable = "users";
      if (urlParams.get("page") == null) {
        history.pushState({}, null, "/moderation&table=users");
      } else {
        history.pushState(
          {},
          null,
          "/moderation?page=" + urlParams.get("page") + "&table=users"
        );
      }
    }
  }
  if (!categoryTable.className.includes("Invisible")) {
    categoryTable.classList.add("Invisible");
    categoryTable.classList.remove("Visible");
  }
}
//Affiche seulement le tableau post et cache les autres. Chaque fonction est casi pareil, je cache les autres table SI elle sont visible et j'affiche celle si si elle ne l'est pas.
function showCategory() {
  if (!commentTable.className.includes("Invisible")) {
    commentTable.classList.add("Invisible");
    commentTable.classList.remove("Visible");
  }
  if (!postTable.className.includes("Invisible")) {
    postTable.classList.add("Invisible");
    postTable.classList.remove("Visible");
  }
  if (!userTable.className.includes("Invisible")) {
    userTable.classList.add("Invisible");
    userTable.classList.remove("Visible");
  }
  if (categoryTable.className.includes("Invisible")) {
    removeSelected();
    removeVisibility("postOption");
    categoryTable.classList.remove("Invisible");
    categoryTable.classList.add("Visible");
    if (actualTable != "categories") {
      actualTable = "categories";
      if (urlParams.get("page") == null) {
        history.pushState({}, null, "/moderation&table=categories");
      } else {
        history.pushState(
          {},
          null,
          "/moderation?page=" + urlParams.get("page") + "&table=categories"
        );
      }
    }
  }
}*/
