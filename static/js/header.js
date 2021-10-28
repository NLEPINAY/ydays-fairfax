class Search {
    constructor(config) {
        this.config = config

        this.padding = 10
        this.deploy = false

        this.containerSearch = document.querySelector(this.config.container)
        this.searchbar = document.querySelector(this.config.searchbar)
        this.iconSearch = document.querySelector(this.config.iconSearch)
        this.closeSearch = document.querySelector(this.config.closeSearch)
        this.btnSearch = document.querySelector(this.config.btnSearch)
    }

    init() {
        this.iconSearch.addEventListener('click', () => {this.deployMenu(this.containerSearch, this.closeSearch)})
        this.closeSearch.addEventListener('click', () => {this.undeployMenu(this.containerSearch, this.closeSearch)})

        this.startSearch()
    }
    deployMenu(element, close) {
        element.style = "transform: translate(calc(-100vw + 30px + 25px), -8vh)"
        close.style = "opacity: 1; z-index: 5; visibility: visible;"
        this.deploy = true
    }
    undeployMenu(element, close) {
        element.style = "transform: translate(calc(-100vw + 30px + 25px), calc(-100vh - 50px))"
        close.style = "opacity: 0; z-index: -1; visibility: hidden;"
        this.deploy = false
    }
    startSearch() {
        document.addEventListener('keyup', ev => {
            if (this.searchbar.value != "" && this.deploy == true) {
                if (ev.key == 'Enter') {
                    this.searchUser()
                }
            }
        })
        this.btnSearch.addEventListener('click', () => {
            if (this.searchbar.value != "" && this.deploy == true) {
                this.searchUser()
            }
        })
    }
    searchUser() {
        console.log('start search')
        // if (this.searchbar.value != " " || this.searchbar.value != "") {
        //     let data = new FormData(),
        //     xhr = XMLHttpRequest

        //     data.append('name', this.searchbar.value)
        //     xhr.open('POST', 'searchUser', true)
        //     xhr.send(data)

        //     // xhr.onreadystatechange = function () {
        //     //     if (xhr.readyState == 4 && xhr.status == 200) {
        //     //         window.location.href = Routing.generate('admin_contact_mails')
        //     //     }
        //     // }
        // }
    }
}

const search = new Search({
    container: '.deploy_search',
    searchbar: '.user_search',
    iconSearch: '.search_icon',
    closeSearch: '.close_search',
    btnSearch: '.btn_search'
})
search.init()