<script context="module">
	/**
	 * @type {import('@sveltejs/kit').Load}
	 */
	export async function load({ page, fetch }) {
		const res = await fetch(`${import.meta.env.VITE_API_URL}/room`);

		if (res.ok) {
			return {
				props: {
					rooms: await res.json()
				}
			};
		}

		return {
			status: res.status,
			error: new Error(`Could not load rooms`)
		};
	}
</script>

<script lang="ts">
	import type { SimpleRoom } from '$lib/types/Room';

	export let rooms: SimpleRoom[];
</script>

<div class="container mx-auto max-w-lg py-4">
	<div class="card shadow-lg">
		<div class="card-body">
			<div class="card-title">Rooms</div>
			<table class="table w-full">
				<thead>
					<tr>
						<th>name</th>
						<th>question</th>
					</tr>
				</thead>
				<tbody>
					{#each rooms as room}
						<tr class="hover">
							<td>{room.id}</td>
							<td>{room.question}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
