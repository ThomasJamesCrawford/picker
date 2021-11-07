<script>
	import '../app.css';
</script>

<svelte:head>
	<title>Pick Now</title>
</svelte:head>

<div class="navbar mb-2 shadow-lg bg-neutral text-neutral-content rounded-box">
	<div class="flex-1 px-2 mx-2">
		<span class="text-lg font-bold">Pick Now</span>
	</div>
	<div class="flex-none hidden px-2 mx-2 lg:flex">
		<div class="flex items-stretch">
			<a href="/" class="btn btn-ghost btn-sm rounded-btn">new room </a>
			<a href="/admin" class="btn btn-ghost btn-sm rounded-btn">my rooms </a>
		</div>
	</div>
</div>

<div class="bg-gray-50 min-h-screen">
	<slot />
</div>
