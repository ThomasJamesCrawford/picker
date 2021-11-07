<script context="module">
	export function load({ error, status }) {
		return {
			props: {
				error,
				status
			}
		};
	}
</script>

<script lang="ts">
	export let error: string;
	export let status: string;
</script>

<div class="contaier mx-auto max-w-lg p-16">
	<div class="mb-8 text-center">{status}</div>
	<div class="alert alert-error">
		<div class="flex-1 justify-center">
			<svg
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				class="w-6 h-6 mx-2 stroke-current"
			>
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
				/>
			</svg>
			<!-- svelte-ignore a11y-label-has-associated-control -->
			<label>{error}</label>
		</div>
	</div>
</div>
