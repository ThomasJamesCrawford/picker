<script>
	import { APP_NAME } from '$lib/constants';

	import '../app.css';
</script>

<svelte:head>
	<title>{APP_NAME}</title>
</svelte:head>

<slot />
