<script lang="ts">
    	import type { LayoutData } from './$types';	
</script>

<main>
<slot/>
</main>

<style lang="scss">
	main {
		height: 100vh;
		display: flex;
		justify-content: center;
		align-items: center;
	}
</style>