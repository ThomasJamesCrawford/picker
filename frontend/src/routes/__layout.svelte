<script>
	import { APP_NAME } from '$lib/constants';

	import '../app.css';
</script>

<svelte:head>
	<title>{APP_NAME}</title>
</svelte:head>

<div class="bg-gray-50 min-h-full">
	<slot />
</div>
