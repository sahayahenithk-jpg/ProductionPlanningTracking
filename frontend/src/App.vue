<template>
  <div class="app-shell">
    <header>
      <h1>Project Planning Tracking</h1>
    </header>
    <main>
      <router-view />
    </main>
  </div>
</template>

<style>
body {
  margin: 0;
  font-family: system-ui, sans-serif;
  background: #f7fafc;
}

.app-shell {
  max-width: 520px;
  margin: 0 auto;
  padding: 32px 16px;
}

header {
  text-align: center;
  margin-bottom: 24px;
}

h1 {
  margin: 0;
  font-size: 28px;
}
</style>
