<div class="row text-center mb-5">
	<div class="col-sm-12">
		<h2>RDP Control Panel</h2>
	</div>
	{if $rdpControlPanelURL}
		<div class="col-sm-12">
			<p>Click the button below to continue</p>
			<a class="btn btn-primary" href="{$rdpControlPanelURL}" target="_blank" rel="noopener noreferrer">Login to RDP Control Panel</a>
		</div>
	{else}
		<div class="col-sm-12" style="cursor: not-allowed;">
			<p>Single Sign-On is not enabled</p>
			<a class="btn btn-primary disabled">Login to RDP Control Panel</a>
		</div>
	{/if}
</div>
