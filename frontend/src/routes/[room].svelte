<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch }) {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/publicRoom/${page.params['room']}`);

		if (res.ok) {
			return {
				props: {
					room: await res.json()
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load that room`)
		};
	}
</script>

<script lang="ts">
	import type { PublicRoom } from '$lib/types/Room';

	export let room: PublicRoom;
</script>

<div class="card shadow-lg">
	<div class="card-title">
		{room.question}
	</div>
	<div class="card-body">
		{#each room.options as option}
			<button type="button" class="btn btn-secondary">
				{option}
			</button>
		{/each}
	</div>
</div>
