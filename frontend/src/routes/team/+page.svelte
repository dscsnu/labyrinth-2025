<script>
    let name = 'world';
  </script>
  
  <main>
    <h1>Hello {name}!</h1>
    <input bind:value={name} placeholder="Enter your name" />
  </main>
  
  <style>
    main {
      padding: 2rem;
      font-family: sans-serif;
      text-align: center;
    }
  </style>
  