Vue.component('owncast-footer', {
  template: `
    <footer class="flex border-t border-gray-500 border-solid">
      <span>
        <a href="https://github.com/gabek/owncast" target="_blank">About Owncast</a>
      </span>          
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