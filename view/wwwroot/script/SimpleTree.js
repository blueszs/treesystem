/*
Date：2011.11.27 0:33
Version：SimpleTree 1.1
*/

$(function(){
	$.fn.extend({
		SimpleTree:function(options){
			
			//初始化参数
			var option = $.extend({
				click:function(a){ }
			},options);
			
			option.tree=this;	/* 在参数对象中添加对当前菜单树的引用，以便在对象中使用该菜单树 */
			
			option._init=function(){
				/*
				 * 初始化菜单展开状态，以及分叉节点的样式
				 */				
				this.tree.find("ul ul").hide();	/* 隐藏所有子级菜单 */
				this.tree.find("ul ul").prev("li").removeClass("open");	/* 移除所有子级菜单父节点的 open 样式 */
				
				this.tree.find("ul ul[show='true']").show();	/* 显示 show 属性为 true 的子级菜单 */
				this.tree.find("ul ul[show='true']").prev("li").addClass("open");	/* 添加 show 属性为 true 的子级菜单父节点的 open 样式 */
			}/* option._init() End */
			
			/* 设置所有超链接不响应单击事件 */
			//this.find("a").click(function(){ $(this).parents("li").click(); return false; });
			
			/* 菜单项 <li> 接受单击 */
			this.find("a").click(function(){
				/*
				 * 当单击菜单项 <li>
				 * 1.触发用户自定义的单击事件，将该 <li> 标签中的第一个超链接做为参数传递过去
				 * 2.修改当前菜单项所属的子菜单的显示状态（如果等于 true 将其设置为 false，否则将其设置为 true）
				 * 3.重新初始化菜单
				 */
				var a=$(this).find("a")[0];
				if(typeof(a)!="undefined")
					option.click(a);	/* 如果获取的超链接不是 undefined，则触发单击 */
				
				/* 
				 * 如果当前节点下面包含子菜单，并且其 show 属性的值为 true，则修改其 show 属性为 false
				 * 否则修改其 show 属性为 true
				 */
				if ($(this).siblings("ul").attr("show") == "true") {
					$(this).siblings("ul").attr("show","false");					
				}else{
					$(this).siblings("ul").attr("show","true");
				}
				
				/* 初始化菜单 */
				option._init();
			});
			
			this.find("li").hover(
				function(){
					$(this).addClass("hover");
				},
				function(){
					$(this).removeClass("hover");
				}
			);
			
			/* 设置所有父节点样式 */
			this.find("ul").prev("li").addClass("folder");
			
			/* 设置节点“是否包含子节点”属性 */
			this.find("li").find("a").attr("hasChild",false);
			this.find("ul").prev("li").find("a").attr("hasChild",true);
			
			/* 初始化菜单 */
			option._init();
			
		}/* SimpleTree Function End */
		
	});
});