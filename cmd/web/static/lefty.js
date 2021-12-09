var app = new Vue({
    el: '#app',
    data: {
        "products": [],
    },
    mounted: async function () {
        const response = await fetch("/api/products");
        const json = await response.json();
        this.products = json.data;
        console.log(this.products);
    }
})