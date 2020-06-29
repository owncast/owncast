Vue.component('owncast-footer', {
  props: {
    appVersion: {
      type: String,
      default: '0.1',
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

Vue.component('user-details', {
  props:  ['logo', 'platforms', 'summary', 'tags', 'extraContent'],
  template: `
    <div class="user-content">
      <div class="user-image rounded-full bg-white">
        <img class="logo" alt="Logo" v-bind:src="logo">            
      </div>
      <div class="user-content-header border-b border-gray-500 border-solid">
        <h2 class="font-semibold">
          About <span class="streamer-name text-indigo-600">
            <slot></slot>
          </span>
        </h2>
        <social-list v-bind:platforms="platforms"></social-list>
        <div class="stream-summary" v-html="summary"></div>
        <stream-tags v-bind:tags="tags"></stream-tags>
      </div>

      <div v-html="extraContent"></div>

    </div>
  `,
});
