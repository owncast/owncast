Vue.component('owncast-footer', {
  props: {
    appVersion: {
      type: String,
      default: "0.1",
    },
  },
  
  template: `
    <footer class="flex">
      <span>
        <a href="https://github.com/gabek/owncast" target="_blank">About Owncast</a>
      </span>
      <span>Version {{appVersion}}</span>
    </footer>
  `,
});


Vue.component('stream-tags', {
  props: ['tags'],
  template: `
    <ul
      class="tag-list flex"
      v-if="this.tags.length"
    >
      <li class="tag rounded-sm text-gray-100 bg-gray-700" 
        v-for="tag in this.tags"
        v-bind:key="tag"
      >
        {{tag}}
      </li>
    </ul>
  `,
});