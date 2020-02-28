<!DOCTYPE html>
<html lang="en">
<head>
	<link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.11.2/css/all.min.css">
    <link href="https://fonts.googleapis.com/css?family=Lato" rel="stylesheet">
    <link rel="stylesheet" href="style.css">
</head>
<body>
	<div class="main-section">
		<p style="color:#fff">
		Last updated at {{.Now.Format "02 Jan 06 15:04 MST"}}. See <a href="https://codefresh.io/steps/">https://codefresh.io/steps/</a> for more details. Next update in 1 hour.
		</p>

		
	{{range $item, $step := .FinishedSteps}}

           	<div class="dashbord {{if (eq $step.Status 1)}} valid-step {{else if (eq $step.Status 2)}} invalid-step {{end}}">
			<div class="title-section">
				<p>{{$step.Status}} {{$step.Name}} {{$step.Version}}</p>
			</div>
			<div class="icon-text-section">
				<div class="icon-section">
					<i class="{{if (eq $step.Status 1)}} fa fa-check {{else if (eq $step.Status 2)}} fa fa-frown {{else}} fab fa-docker {{end}}" aria-hidden="true"></i>
				</div>
				<div class="text-section">
					<h1>{{len $step.ImagesUsed }} Docker image(s) used</h1>
					<span>+7% email list penetration</span>
				</div>
				<div style="clear:both;"></div>
			</div>
			<div class="detail-section">
				<a href="{{$step.SourceURL}}">
					<p>View Source code</p>
					<i class="fa fa-arrow-right" aria-hidden="true"></i>
				</a>
			</div>
		</div>

    {{end}}





	</div>

</body>
</html>
