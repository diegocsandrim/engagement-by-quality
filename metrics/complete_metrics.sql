select
	projects.kee as project_name,
	to_timestamp(snapshots.created_at/1000)::date date,
	substring(substring(convert_from(ce_scanner_context.context_data, 'UTF-8'), 'sonar.analysis.contributorGain=\d*') from length('sonar.analysis.contributorGain=')+1)::integer as contributor_gain,
	metrics.name as metric_name,
	replace(project_measures.value::text, '.', ',') AS metric_value
from projects 
	left join ce_activity
		on projects.project_uuid=ce_activity.component_uuid
	left join snapshots
		on snapshots.uuid = ce_activity.analysis_uuid
	left join ce_scanner_context
		on ce_activity.uuid = ce_scanner_context.task_uuid
	left join project_measures
		on ce_activity.analysis_uuid = project_measures.analysis_uuid
	left join metrics
		on project_measures.metric_id=metrics.id
where
	metrics.name in (
		'blocker_violations',
		'bugs',
		'classes',
		'code_smells',
		'cognitive_complexity',
		'comment_lines',
		'comment_lines_density',
		'complexity',
		'critical_violations',
		'duplicated_blocks',
		'duplicated_files',
		'duplicated_lines',
		'duplicated_lines_density',
		'effort_to_reach_maintainability_rating_a',
		'file_complexity',
		'files',
		'functions',
		'info_violations',
		'lines',
		'major_violations',
		'minor_violations',
		'ncloc',
		'open_issues',
		'reliability_rating',
		'reliability_remediation_effort',
		'security_rating',
		'security_remediation_effort',
		'sqale_debt_ratio',
		'sqale_index',
		'sqale_rating',
		'statements',
		'violations',
		'vulnerabilities'
		)
	and projects.scope='PRJ'
order by
	projects.kee,
	date,
	metrics.name